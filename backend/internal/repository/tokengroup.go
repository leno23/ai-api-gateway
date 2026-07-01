package repository

import "github.com/leno23/ai-api-gateway/internal/model"

func (r *Repos) ListTokenGroups() ([]model.TokenGroup, error) {
	var rows []model.TokenGroup
	if err := r.DB.Where("status = ?", 1).Order("id asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repos) TokenGroupBySlug(slug string) (*model.TokenGroup, error) {
	var g model.TokenGroup
	if err := r.DB.Where("slug = ? AND status = ?", slug, 1).First(&g).Error; err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *Repos) TokenGroupForUser(userID int64) (*model.TokenGroup, error) {
	var u model.User
	if err := r.DB.Select("token_group_id").First(&u, userID).Error; err != nil {
		return nil, err
	}
	if u.TokenGroupID != nil {
		var g model.TokenGroup
		if err := r.DB.First(&g, *u.TokenGroupID).Error; err == nil {
			return &g, nil
		}
	}
	return r.TokenGroupBySlug("default")
}
