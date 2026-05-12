package rebate

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/leno23/ai-api-gateway/internal/model"
)

// EnqueueConsumeRebate schedules a rebate as a fraction (basis points) of consumed quota.
// Default-off: caller passes enabled=false from config.
func EnqueueConsumeRebate(db *gorm.DB, log *zap.Logger, userID, consumedQuota int64, ref string, enabled bool, bps int) {
	if !enabled || bps <= 0 || consumedQuota <= 0 || db == nil {
		return
	}
	amt := consumedQuota * int64(bps) / 10000
	if amt <= 0 {
		return
	}
	go func() {
		defer func() { _ = recover() }()
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		t := model.RebateTask{
			UserID: userID,
			Amount: amt,
			Status: model.RebateTaskPending,
			Ref:    truncate(ref, 128),
		}
		if err := db.WithContext(ctx).Create(&t).Error; err != nil && log != nil {
			log.Warn("rebate enqueue", zap.Error(err), zap.Int64("user_id", userID))
		}
	}()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

// RefChat builds a short reference for chat consumption.
func RefChat(modelName string) string {
	return fmt.Sprintf("chat:%s", modelName)
}
