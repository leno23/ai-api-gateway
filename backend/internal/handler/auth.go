package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/leno23/ai-api-gateway/internal/auth"
	"github.com/leno23/ai-api-gateway/internal/middleware"
	"github.com/leno23/ai-api-gateway/internal/model"
	"github.com/leno23/ai-api-gateway/internal/service"
)

const (
	inviteeBonusTokens = int64(1000)
	inviterBonusTokens = int64(1000)
	jwtTTL             = 72 * time.Hour
)

type AuthDeps struct {
	DB        *gorm.DB
	JWTSecret []byte
}

type registerReq struct {
	Email      string `json:"email" binding:"required,email"`
	Username   string `json:"username" binding:"required,min=3,max=64"`
	Password   string `json:"password" binding:"required,min=8,max=128"`
	InviteCode string `json:"invite_code"`
}

func Register(deps AuthDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req registerReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "hash failed"})
			return
		}
		invite, err := service.RandomInviteCode(8)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invite code"})
			return
		}
		var inviter *model.User
		if req.InviteCode != "" {
			var u model.User
			if err := deps.DB.Where("invite_code = ?", req.InviteCode).First(&u).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invite code"})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": "lookup failed"})
				return
			}
			inviter = &u
		}
		u := model.User{
			Email:        req.Email,
			Username:     req.Username,
			PasswordHash: string(hash),
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
			InviteCode:   invite,
		}
		if inviter != nil {
			id := inviter.ID
			u.InvitedBy = &id
			u.Quota = inviteeBonusTokens
		}
		err = deps.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&u).Error; err != nil {
				return err
			}
			if inviter != nil {
				if err := tx.Model(&model.User{}).Where("id = ?", inviter.ID).
					UpdateColumn("quota", gorm.Expr("quota + ?", inviterBonusTokens)).Error; err != nil {
					return err
				}
				ir := model.InviteRecord{
					InviterID:    inviter.ID,
					InviteeID:    u.ID,
					InviterBonus: inviterBonusTokens,
					InviteeBonus: inviteeBonusTokens,
					Status:       1,
				}
				if err := tx.Create(&ir).Error; err != nil {
					return err
				}
				var inviterBal int64
				if err := tx.Model(&model.User{}).Select("quota").Where("id = ?", inviter.ID).Scan(&inviterBal).Error; err != nil {
					return err
				}
				if err := tx.Create(&model.QuotaLog{
					UserID:    inviter.ID,
					Delta:     inviterBonusTokens,
					Balance:   inviterBal,
					Type:      model.QuotaLogTypeInvite,
					Reference: "",
					Remark:    "invite reward",
				}).Error; err != nil {
					return err
				}
				if err := tx.Create(&model.QuotaLog{
					UserID:    u.ID,
					Delta:     inviteeBonusTokens,
					Balance:   u.Quota,
					Type:      model.QuotaLogTypeInvite,
					Reference: "",
					Remark:    "signup with invite",
				}).Error; err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "email or username taken"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"id":          u.ID,
			"email":       u.Email,
			"username":    u.Username,
			"invite_code": u.InviteCode,
		})
	}
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(deps AuthDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var u model.User
		if err := deps.DB.Where("email = ?", req.Email).First(&u).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		if u.Status != model.UserStatusActive {
			c.JSON(http.StatusForbidden, gin.H{"error": "account disabled"})
			return
		}
		tok, err := auth.SignJWT(deps.JWTSecret, u.ID, u.Role, jwtTTL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"access_token": tok, "token_type": "Bearer", "expires_in": int(jwtTTL.Seconds())})
	}
}

type createKeyReq struct {
	Name string `json:"name" binding:"required,min=1,max=128"`
}

type APIKeyDeps struct {
	DB *gorm.DB
}

func CreateAPIKey(deps APIKeyDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := c.Get(middleware.CtxUserID)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID := uid.(int64)
		var req createKeyReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		raw, err := service.RandomRedeemCode()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rng"})
			return
		}
		plain := "sk-" + raw
		hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "hash"})
			return
		}
		prefix := plain
		if len(prefix) > 16 {
			prefix = prefix[:16]
		}
		k := model.APIKey{
			UserID:    userID,
			Name:      req.Name,
			KeyHash:   string(hash),
			KeyPrefix: prefix,
			Status:    model.APIKeyStatusActive,
		}
		if err := deps.DB.Create(&k).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "persist"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"id":         k.ID,
			"name":       k.Name,
			"key_prefix": k.KeyPrefix,
			"api_key":    plain,
		})
	}
}
