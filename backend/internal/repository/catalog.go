package repository

import (
	"strings"

	"github.com/leno23/ai-api-gateway/internal/model"
)

type CatalogFilter struct {
	Provider     string
	EndpointType string
	BillingType  *int16
	Tag          string
	Query        string
	Page         int
	PageSize     int
}

func (r *Repos) ListCatalogModels(f CatalogFilter) ([]model.ModelPrice, int64, error) {
	q := r.DB.Model(&model.ModelPrice{})
	if f.Provider != "" {
		q = q.Where("provider = ?", f.Provider)
	}
	if f.EndpointType != "" {
		q = q.Where("endpoint_type = ?", f.EndpointType)
	}
	if f.BillingType != nil {
		q = q.Where("billing_type = ?", *f.BillingType)
	}
	if f.Tag != "" {
		q = q.Where("? = ANY(tags)", f.Tag)
	}
	if f.Query != "" {
		like := "%" + strings.ToLower(f.Query) + "%"
		q = q.Where(
			"LOWER(model) LIKE ? OR LOWER(COALESCE(display_name, '')) LIKE ?",
			like, like,
		)
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
		size = 20
	}
	if size > 200 {
		size = 200
	}
	var rows []model.ModelPrice
	if err := q.Order("provider asc, model asc").
		Offset((page - 1) * size).
		Limit(size).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *Repos) ModelPriceByName(modelName string) (*model.ModelPrice, error) {
	var mp model.ModelPrice
	if err := r.DB.Where("model = ?", modelName).First(&mp).Error; err != nil {
		return nil, err
	}
	return &mp, nil
}

func (r *Repos) DistinctProviders() ([]string, error) {
	var providers []string
	err := r.DB.Model(&model.ModelPrice{}).
		Distinct("provider").
		Order("provider asc").
		Pluck("provider", &providers).Error
	return providers, err
}
