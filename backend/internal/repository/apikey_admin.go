package repository

import (
	"errors"

	"github.com/leno23/ai-api-gateway/internal/model"
	"gorm.io/gorm"
)

var ErrAPIKeyNotFound = errors.New("api key not found")

func (r *Repos) ListAPIKeysByUser(userID int64, page, pageSize int) ([]model.APIKey, int64, error) {
	q := r.DB.Model(&model.APIKey{}).Where("user_id = ?", userID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	var rows []model.APIKey
	if err := q.Order("id desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *Repos) GetAPIKeyForUser(userID, keyID int64) (*model.APIKey, error) {
	var k model.APIKey
	if err := r.DB.Where("id = ? AND user_id = ?", keyID, userID).First(&k).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAPIKeyNotFound
		}
		return nil, err
	}
	return &k, nil
}

func (r *Repos) DeleteAPIKeysForUser(userID int64, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	res := r.DB.Where("user_id = ? AND id IN ?", userID, ids).Delete(&model.APIKey{})
	return res.RowsAffected, res.Error
}
