package circuit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	failKeyTTL = 20 * time.Second
	openKeyTTL = 45 * time.Second
	failOpenAt = int64(5)
)

func failKey(chID int64) string {
	return fmt.Sprintf("cb:ch:%d:fail", chID)
}

func openKey(chID int64) string {
	return fmt.Sprintf("cb:ch:%d:open", chID)
}

// Allow reports whether the channel is not in open (trip) state.
func Allow(ctx context.Context, rdb *redis.Client, chID int64) bool {
	if rdb == nil {
		return true
	}
	v, err := rdb.Get(ctx, openKey(chID)).Result()
	if err == redis.Nil {
		return true
	}
	if err != nil {
		return true
	}
	return v == ""
}

// RecordFailure increments recent failures; opens circuit when threshold is hit.
func RecordFailure(ctx context.Context, rdb *redis.Client, chID int64) {
	if rdb == nil {
		return
	}
	n, err := rdb.Incr(ctx, failKey(chID)).Result()
	if err != nil {
		return
	}
	_ = rdb.Expire(ctx, failKey(chID), failKeyTTL).Err()
	if n >= failOpenAt {
		_ = rdb.Set(ctx, openKey(chID), "1", openKeyTTL).Err()
		_ = rdb.Del(ctx, failKey(chID)).Err()
	}
}

// RecordSuccess clears failure streak for a channel.
func RecordSuccess(ctx context.Context, rdb *redis.Client, chID int64) {
	if rdb == nil {
		return
	}
	_ = rdb.Del(ctx, failKey(chID)).Err()
}
