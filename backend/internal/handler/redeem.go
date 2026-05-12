package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/leno23/ai-api-gateway/internal/middleware"
	"github.com/leno23/ai-api-gateway/internal/model"
	"github.com/leno23/ai-api-gateway/internal/service"
)

type redeemReq struct {
	Code string `json:"code" binding:"required"`
}

type RedeemDeps struct {
	DB *gorm.DB
}

func RedeemUser(deps RedeemDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := c.Get(middleware.CtxUserID)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID := uid.(int64)
		var req redeemReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var outQuota int64
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
		c.JSON(http.StatusOK, gin.H{"quota": outQuota})
	}
}

type batchRedeemReq struct {
	Count       int   `json:"count" binding:"required,min=1,max=5000"`
	Quota       int64 `json:"quota" binding:"required,min=1"`
	ExpiresDays *int  `json:"expires_days"`
}

type AdminDeps struct {
	DB *gorm.DB
}

func BatchRedeem(deps AdminDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req batchRedeemReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var exp *time.Time
		if req.ExpiresDays != nil && *req.ExpiresDays > 0 {
			t := time.Now().Add(time.Duration(*req.ExpiresDays) * 24 * time.Hour)
			exp = &t
		}
		codes := make([]string, 0, req.Count)
		for i := 0; i < req.Count; i++ {
			var code string
			for attempt := 0; attempt < 5; attempt++ {
				raw, err := service.RandomRedeemCode()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "rng"})
					return
				}
				rc := model.RedeemCode{
					Code:      raw,
					Quota:     req.Quota,
					Status:    model.RedeemStatusAvailable,
					ExpiresAt: exp,
				}
				if err := deps.DB.Create(&rc).Error; err == nil {
					code = raw
					break
				}
			}
			if code == "" {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not allocate unique code"})
				return
			}
			codes = append(codes, code)
		}
		c.JSON(http.StatusCreated, gin.H{"codes": codes, "count": len(codes)})
	}
}
