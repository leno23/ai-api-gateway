package jobs

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/leno23/ai-api-gateway/internal/model"
	"github.com/leno23/ai-api-gateway/internal/service"
)

// StartQuotaReconcile periodically overwrites Redis quota:* keys from Postgres (heals cache drift).
func StartQuotaReconcile(ctx context.Context, db *gorm.DB, rdb *redis.Client, log *zap.Logger, interval time.Duration) {
	if interval <= 0 || rdb == nil || db == nil {
		return
	}
	t := time.NewTicker(interval)
	go func() {
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if err := reconcileQuotaOnce(context.Background(), db, rdb, log); err != nil && log != nil {
					log.Warn("quota reconcile", zap.Error(err))
				}
			}
		}
	}()
}

func reconcileQuotaOnce(ctx context.Context, db *gorm.DB, rdb *redis.Client, log *zap.Logger) error {
	var cur uint64
	prefix := "quota:user:"
	for {
		keys, next, err := rdb.Scan(ctx, cur, prefix+"*", 64).Result()
		if err != nil {
			return err
		}
		for _, key := range keys {
			idStr := strings.TrimPrefix(key, prefix)
			uid, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				continue
			}
			var q int64
			if err := db.WithContext(ctx).Model(&model.User{}).Select("quota").Where("id = ?", uid).Scan(&q).Error; err != nil {
				continue
			}
			if err := service.SyncQuotaToRedis(ctx, rdb, uid, q); err != nil && log != nil {
				log.Warn("quota reconcile sync", zap.Int64("user_id", uid), zap.Error(err))
			}
		}
		cur = next
		if cur == 0 {
			break
		}
	}
	return nil
}
