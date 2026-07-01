package repository

import (
	"errors"
	"net"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/leno23/ai-api-gateway/internal/model"
)

var (
	ErrInvalidAPIKey   = errors.New("invalid api key")
	ErrAPIKeyLookup    = errors.New("api key lookup failed")
	ErrUserInactive    = errors.New("user inactive")
	ErrAPIKeyDisabled  = errors.New("api key disabled")
	ErrAPIKeyIPDenied  = errors.New("ip not allowed")
	ErrAPIKeyQuota     = errors.New("api key quota exceeded")
)

// GatewayKeyAuth is the resolved gateway API key context.
type GatewayKeyAuth struct {
	UserID      int64
	KeyID       int64
	Models      []string
	QuotaLimit  *int64
	UsedQuota   int64
	IPWhitelist []string
}

// AuthenticateGatewayAPIKey resolves `sk-...` to an active user and key policy.
func (r *Repos) AuthenticateGatewayAPIKey(raw string, prefixLen int, clientIP string) (*GatewayKeyAuth, error) {
	if len(raw) < prefixLen {
		return nil, ErrInvalidAPIKey
	}
	prefix := raw[:prefixLen]
	var keys []model.APIKey
	if err := r.DB.Where("key_prefix = ?", prefix).Find(&keys).Error; err != nil {
		return nil, ErrAPIKeyLookup
	}
	var matched *model.APIKey
	for i := range keys {
		if err := bcrypt.CompareHashAndPassword([]byte(keys[i].KeyHash), []byte(raw)); err == nil {
			matched = &keys[i]
			break
		}
	}
	if matched == nil {
		return nil, ErrInvalidAPIKey
	}
	if matched.Status != model.APIKeyStatusActive {
		return nil, ErrAPIKeyDisabled
	}
	var u model.User
	if err := r.DB.First(&u, matched.UserID).Error; err != nil || u.Status != model.UserStatusActive {
		return nil, ErrUserInactive
	}
	whitelist := []string(matched.IPWhitelist)
	if len(whitelist) > 0 && !ipAllowed(clientIP, whitelist) {
		return nil, ErrAPIKeyIPDenied
	}
	if matched.QuotaLimit != nil && matched.UsedQuota >= *matched.QuotaLimit {
		return nil, ErrAPIKeyQuota
	}
	return &GatewayKeyAuth{
		UserID:      u.ID,
		KeyID:       matched.ID,
		Models:      []string(matched.Models),
		QuotaLimit:  matched.QuotaLimit,
		UsedQuota:   matched.UsedQuota,
		IPWhitelist: whitelist,
	}, nil
}

func ipAllowed(clientIP string, whitelist []string) bool {
	if clientIP == "" {
		return false
	}
	ip := net.ParseIP(clientIP)
	if ip == nil {
		return false
	}
	for _, entry := range whitelist {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if _, cidr, err := net.ParseCIDR(entry); err == nil {
			if cidr.Contains(ip) {
				return true
			}
			continue
		}
		if host := net.ParseIP(entry); host != nil && host.Equal(ip) {
			return true
		}
	}
	return false
}
