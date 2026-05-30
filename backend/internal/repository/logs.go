package repository

import (
	"strings"
	"time"

	"github.com/leno23/ai-api-gateway/internal/model"
)

type UsageLogFilter struct {
	UserID     int64
	Start      *time.Time
	End        *time.Time
	TokenName  string
	Model      string
	RequestID  string
	TokenGroup string
	Page       int
	PageSize   int
}

func (r *Repos) ListUsageLogs(f UsageLogFilter) ([]model.RequestLog, int64, error) {
	q := r.DB.Model(&model.RequestLog{}).Where("user_id = ?", f.UserID)
	if f.Start != nil {
		q = q.Where("created_at >= ?", *f.Start)
	}
	if f.End != nil {
		q = q.Where("created_at <= ?", *f.End)
	}
	if f.TokenName != "" {
		q = q.Where("token_name ILIKE ?", "%"+f.TokenName+"%")
	}
	if f.Model != "" {
		q = q.Where("model ILIKE ?", "%"+f.Model+"%")
	}
	if f.RequestID != "" {
		q = q.Where("request_id = ?", f.RequestID)
	}
	if f.TokenGroup != "" {
		q = q.Where("token_group = ?", f.TokenGroup)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page := f.Page
	if page < 1 {
		page = 1
	}
	size := f.PageSize
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	var rows []model.RequestLog
	if err := q.Order("created_at desc").
		Offset((page - 1) * size).
		Limit(size).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

type TaskLogFilter struct {
	UserID  int64
	Start   *time.Time
	End     *time.Time
	TaskID  string
	Page    int
	PageSize int
}

func (r *Repos) ListTaskLogs(f TaskLogFilter) ([]model.TaskLog, int64, error) {
	q := r.DB.Model(&model.TaskLog{}).Where("user_id = ?", f.UserID)
	if f.Start != nil {
		q = q.Where("submitted_at >= ?", *f.Start)
	}
	if f.End != nil {
		q = q.Where("submitted_at <= ?", *f.End)
	}
	if f.TaskID != "" {
		q = q.Where("task_id ILIKE ?", "%"+strings.TrimSpace(f.TaskID)+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page := f.Page
	if page < 1 {
		page = 1
	}
	size := f.PageSize
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	var rows []model.TaskLog
	if err := q.Order("submitted_at desc").
		Offset((page - 1) * size).
		Limit(size).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}
