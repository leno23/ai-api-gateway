package service

import (
	"gorm.io/gorm"

	"github.com/leno23/ai-api-gateway/internal/model"
)

// CreditAffiliateOnRecharge adds a share of recharge quota to the inviter's pending balance.
func CreditAffiliateOnRecharge(tx *gorm.DB, inviteeID int64, rechargeQuota int64, bps int) error {
	if bps <= 0 || rechargeQuota <= 0 {
		return nil
	}
	var invitee model.User
	if err := tx.Select("invited_by").First(&invitee, inviteeID).Error; err != nil {
		return nil
	}
	if invitee.InvitedBy == nil {
		return nil
	}
	bonus := rechargeQuota * int64(bps) / 10000
	if bonus <= 0 {
		return nil
	}
	return tx.Model(&model.User{}).Where("id = ?", *invitee.InvitedBy).
		UpdateColumn("affiliate_pending", gorm.Expr("affiliate_pending + ?", bonus)).Error
}
