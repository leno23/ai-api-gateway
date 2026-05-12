package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const quotaKeyTTL = 24 * time.Hour

func quotaRedisKey(userID int64) string {
	return fmt.Sprintf("quota:user:%d", userID)
}

// SyncQuotaToRedis refreshes cached quota from DB (best-effort).
func SyncQuotaToRedis(ctx context.Context, rdb *redis.Client, userID, quota int64) error {
	return rdb.Set(ctx, quotaRedisKey(userID), strconv.FormatInt(quota, 10), quotaKeyTTL).Err()
}

// PreDeduct seeds Redis from DB when missing, then decrements by cost. Returns balance after deduct.
func PreDeduct(ctx context.Context, rdb *redis.Client, db *gorm.DB, userID int64, cost int64) (int64, error) {
	key := quotaRedisKey(userID)
	n, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if n == 0 {
		var u struct {
			Quota int64
		}
		if err := db.Table("users").Select("quota").Where("id = ?", userID).Take(&u).Error; err != nil {
			return 0, err
		}
		if err := rdb.Set(ctx, key, strconv.FormatInt(u.Quota, 10), quotaKeyTTL).Err(); err != nil {
			return 0, err
		}
	}
	bal, err := rdb.DecrBy(ctx, key, cost).Result()
	if err != nil {
		return 0, err
	}
	return bal, nil
}

// RefundQuota increments Redis balance after failed upstream (best-effort).
func RefundQuota(ctx context.Context, rdb *redis.Client, userID, amount int64) error {
	key := quotaRedisKey(userID)
	return rdb.IncrBy(ctx, key, amount).Err()
}
