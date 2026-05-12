package jobs

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/leno23/ai-api-gateway/internal/model"
	"github.com/leno23/ai-api-gateway/internal/service"
)

// StartRebateWorker applies pending rebate_tasks (credits user quota + quota_logs).
func StartRebateWorker(ctx context.Context, db *gorm.DB, rdb *redis.Client, log *zap.Logger, interval time.Duration) {
	if interval <= 0 || db == nil {
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
				if err := processRebateBatch(context.Background(), db, rdb, log); err != nil && log != nil {
					log.Warn("rebate batch", zap.Error(err))
				}
			}
		}
	}()
}

func processRebateBatch(ctx context.Context, db *gorm.DB, rdb *redis.Client, log *zap.Logger) error {
	var ids []int64
	if err := db.WithContext(ctx).Model(&model.RebateTask{}).
		Where("status = ?", model.RebateTaskPending).
		Order("id ASC").
		Limit(80).
		Pluck("id", &ids).Error; err != nil {
		return err
	}
	for _, id := range ids {
		var committedUID int64
		var committedQ int64
		err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var task model.RebateTask
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&task, id).Error; err != nil {
				return err
			}
			if task.Status != model.RebateTaskPending {
				return nil
			}
			if err := tx.Model(&model.User{}).Where("id = ?", task.UserID).
				UpdateColumn("quota", gorm.Expr("quota + ?", task.Amount)).Error; err != nil {
				return err
			}
			var qAfter int64
			if err := tx.Model(&model.User{}).Select("quota").Where("id = ?", task.UserID).Scan(&qAfter).Error; err != nil {
				return err
			}
			now := time.Now()
			if err := tx.Create(&model.QuotaLog{
				UserID:    task.UserID,
				Delta:     task.Amount,
				Balance:   qAfter,
				Type:      model.QuotaLogTypeConsumeRebate,
				Reference: task.Ref,
				Remark:    "consume rebate",
			}).Error; err != nil {
				return err
			}
			committedUID = task.UserID
			committedQ = qAfter
			return tx.Model(&task).Updates(map[string]any{
				"status":       model.RebateTaskDone,
				"processed_at": now,
			}).Error
		})
		if err != nil {
			now := time.Now()
			_ = db.WithContext(ctx).Model(&model.RebateTask{}).Where("id = ?", id).Updates(map[string]any{
				"status":       model.RebateTaskFailed,
				"processed_at": now,
			}).Error
			if log != nil {
				log.Warn("rebate task failed", zap.Error(err), zap.Int64("task_id", id))
			}
			continue
		}
		if committedUID > 0 && rdb != nil {
			_ = service.SyncQuotaToRedis(ctx, rdb, committedUID, committedQ)
		}
	}
	return nil
}
