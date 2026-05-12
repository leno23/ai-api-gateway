package repository

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/leno23/ai-api-gateway/internal/model"
)

var (
	ErrInvalidAPIKey = errors.New("invalid api key")
	ErrAPIKeyLookup  = errors.New("api key lookup failed")
	ErrUserInactive  = errors.New("user inactive")
)

// AuthenticateGatewayAPIKey resolves `sk-...` to an active user and returns allowed models (empty = all).
func (r *Repos) AuthenticateGatewayAPIKey(raw string, prefixLen int) (userID, keyID int64, models []string, err error) {
	if len(raw) < prefixLen {
		return 0, 0, nil, ErrInvalidAPIKey
	}
	prefix := raw[:prefixLen]
	var keys []model.APIKey
	if err := r.DB.Where("key_prefix = ? AND status = ?", prefix, model.APIKeyStatusActive).Find(&keys).Error; err != nil {
		return 0, 0, nil, ErrAPIKeyLookup
	}
	var matched *model.APIKey
	for i := range keys {
		if err := bcrypt.CompareHashAndPassword([]byte(keys[i].KeyHash), []byte(raw)); err == nil {
			matched = &keys[i]
			break
		}
	}
	if matched == nil {
		return 0, 0, nil, ErrInvalidAPIKey
	}
	var u model.User
	if err := r.DB.First(&u, matched.UserID).Error; err != nil || u.Status != model.UserStatusActive {
		return 0, 0, nil, ErrUserInactive
	}
	return u.ID, matched.ID, []string(matched.Models), nil
}
