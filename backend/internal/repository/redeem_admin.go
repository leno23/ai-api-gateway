package repository

import (
	"github.com/leno23/ai-api-gateway/internal/model"
)

// RedeemListFilter filters admin redeem code listing.
type RedeemListFilter struct {
	Status *int16
	Offset int
	Limit  int
}

// ListRedeemCodes returns a page of redeem_codes ordered by id desc.
func (r *Repos) ListRedeemCodes(f RedeemListFilter) ([]model.RedeemCode, int64, error) {
	q := r.DB.Model(&model.RedeemCode{})
	if f.Status != nil {
		q = q.Where("status = ?", *f.Status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.RedeemCode
	if err := q.Order("id DESC").Offset(f.Offset).Limit(f.Limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

type redeemStatRow struct {
	Status int16 `gorm:"column:status"`
	Cnt    int64 `gorm:"column:cnt"`
}

// RedeemCodeStats returns counts grouped by status.
func (r *Repos) RedeemCodeStats() (map[int16]int64, error) {
	var rows []redeemStatRow
	if err := r.DB.Model(&model.RedeemCode{}).
		Select("status, COUNT(*) AS cnt").
		Group("status").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[int16]int64, len(rows))
	for _, row := range rows {
		out[row.Status] = row.Cnt
	}
	return out, nil
}
