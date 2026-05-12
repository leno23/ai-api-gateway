package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func minuteWindow() int64 {
	return time.Now().Unix() / 60
}

// AllowFixedWindow increments a per-minute counter and returns false when above limit.
func AllowFixedWindow(ctx context.Context, rdb *redis.Client, key string, limit int64) (bool, error) {
	if rdb == nil || limit <= 0 {
		return true, nil
	}
	n, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return true, err
	}
	if n == 1 {
		_ = rdb.Expire(ctx, key, 2*time.Minute).Err()
	}
	return n <= limit, nil
}

func GlobalKey() string {
	return fmt.Sprintf("rl:g:%d", minuteWindow())
}

func UserKey(userID int64) string {
	return fmt.Sprintf("rl:u:%d:%d", userID, minuteWindow())
}

func ChannelKey(chID int64) string {
	return fmt.Sprintf("rl:ch:%d:%d", chID, minuteWindow())
}
