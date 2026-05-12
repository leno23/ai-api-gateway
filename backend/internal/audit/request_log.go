package audit

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/leno23/ai-api-gateway/internal/model"
)

// SubmitRequestLog persists a request log asynchronously (best-effort).
// Quota / quota_logs must already be committed synchronously.
func SubmitRequestLog(db *gorm.DB, log *zap.Logger, rec *model.RequestLog) {
	go func() {
		defer func() {
			_ = recover()
		}()
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := db.WithContext(ctx).Create(rec).Error; err != nil {
			if log != nil {
				log.Warn("async request_log", zap.Error(err))
			}
		}
	}()
}
