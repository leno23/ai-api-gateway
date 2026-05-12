package repository

import "github.com/leno23/ai-api-gateway/internal/model"

func (r *Repos) ListModelPrices() ([]model.ModelPrice, error) {
	var rows []model.ModelPrice
	if err := r.DB.Order("model asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
