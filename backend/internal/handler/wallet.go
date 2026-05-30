package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/leno23/ai-api-gateway/internal/model"
	"github.com/leno23/ai-api-gateway/internal/repository"
	"github.com/leno23/ai-api-gateway/internal/service"
)

type WalletDeps struct {
	DB               *gorm.DB
	Repos            *repository.Repos
	RechargeEnabled  bool
	AffiliateBPS     int
	PortalOrigin     string
}

func WalletSummary(deps WalletDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := ctxUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		var u model.User
		if err := deps.DB.First(&u, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		var requestCount int64
		_ = deps.DB.Model(&model.RequestLog{}).Where("user_id = ?", userID).Count(&requestCount).Error

		origin := deps.PortalOrigin
		if origin == "" {
			origin = "http://localhost:3001"
		}
		inviteURL := fmt.Sprintf("%s/register?aff=%s", origin, u.InviteCode)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"quota":              u.Quota,
				"used_quota":         u.UsedQuota,
				"request_count":      requestCount,
				"invite_code":        u.InviteCode,
				"invite_url":         inviteURL,
				"affiliate_pending":  u.AffiliatePending,
				"recharge_enabled":   deps.RechargeEnabled,
			},
		})
	}
}

func RedeemWithAffiliate(deps WalletDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := ctxUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		var req redeemReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var outQuota int64
		var redeemedQuota int64
		err := deps.DB.Transaction(func(tx *gorm.DB) error {
			var rc model.RedeemCode
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("code = ?", req.Code).First(&rc).Error; err != nil {
				return err
			}
			if rc.Status != model.RedeemStatusAvailable {
				return gorm.ErrRecordNotFound
			}
			if rc.ExpiresAt != nil && time.Now().After(*rc.ExpiresAt) {
				return gorm.ErrRecordNotFound
			}
			now := time.Now()
			if err := tx.Model(&model.RedeemCode{}).Where("id = ?", rc.ID).Updates(map[string]any{
				"status":  model.RedeemStatusUsed,
				"used_by": userID,
				"used_at": now,
			}).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.User{}).Where("id = ?", userID).
				UpdateColumn("quota", gorm.Expr("quota + ?", rc.Quota)).Error; err != nil {
				return err
			}
			redeemedQuota = rc.Quota
			if err := service.CreditAffiliateOnRecharge(tx, userID, rc.Quota, deps.AffiliateBPS); err != nil {
				return err
			}
			if err := tx.Model(&model.User{}).Select("quota").Where("id = ?", userID).Scan(&outQuota).Error; err != nil {
				return err
			}
			return tx.Create(&model.QuotaLog{
				UserID:    userID,
				Delta:     rc.Quota,
				Balance:   outQuota,
				Type:      model.QuotaLogTypeRecharge,
				Reference: rc.Code,
				Remark:    "redeem",
			}).Error
		})
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or used code"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "redeem failed"})
			return
		}
		_ = redeemedQuota
		c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"quota": outQuota}})
	}
}

func TransferAffiliate(deps WalletDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := ctxUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		var outQuota int64
		err := deps.DB.Transaction(func(tx *gorm.DB) error {
			var u model.User
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&u, userID).Error; err != nil {
				return err
			}
			if u.AffiliatePending <= 0 {
				return gorm.ErrRecordNotFound
			}
			pending := u.AffiliatePending
			if err := tx.Model(&model.User{}).Where("id = ?", userID).Updates(map[string]any{
				"affiliate_pending": 0,
				"quota":             gorm.Expr("quota + ?", pending),
			}).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.User{}).Select("quota").Where("id = ?", userID).Scan(&outQuota).Error; err != nil {
				return err
			}
			return tx.Create(&model.QuotaLog{
				UserID:    userID,
				Delta:     pending,
				Balance:   outQuota,
				Type:      model.QuotaLogTypeRecharge,
				Reference: "affiliate",
				Remark:    "affiliate transfer",
			}).Error
		})
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusBadRequest, gin.H{"error": "no pending affiliate earnings"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "transfer failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    gin.H{"quota": outQuota, "transferred": true},
		})
	}
}

type mockRechargeReq struct {
	Amount int64 `json:"amount" binding:"required,min=1"`
}

func MockRecharge(deps WalletDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !deps.RechargeEnabled {
			c.JSON(http.StatusForbidden, gin.H{"error": "online recharge disabled"})
			return
		}
		userID, ok := ctxUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		var req mockRechargeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var outQuota int64
		tradeNo := "mock-" + uuid.NewString()
		err := deps.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&model.RechargeRecord{
				UserID:        userID,
				Amount:        float64(req.Amount) / 10000,
				PaymentMethod: "mock",
				TradeNo:       tradeNo,
				Status:        1,
				CompletedAt:   ptrTime(time.Now()),
			}).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.User{}).Where("id = ?", userID).
				UpdateColumn("quota", gorm.Expr("quota + ?", req.Amount)).Error; err != nil {
				return err
			}
			if err := service.CreditAffiliateOnRecharge(tx, userID, req.Amount, deps.AffiliateBPS); err != nil {
				return err
			}
			if err := tx.Model(&model.User{}).Select("quota").Where("id = ?", userID).Scan(&outQuota).Error; err != nil {
				return err
			}
			return tx.Create(&model.QuotaLog{
				UserID:    userID,
				Delta:     req.Amount,
				Balance:   outQuota,
				Type:      model.QuotaLogTypeRecharge,
				Reference: tradeNo,
				Remark:    "mock recharge",
			}).Error
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "recharge failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"quota": outQuota, "trade_no": tradeNo}})
	}
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
