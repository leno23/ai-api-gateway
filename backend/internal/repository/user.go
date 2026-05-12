package repository

import (
	"github.com/leno23/ai-api-gateway/internal/model"
)

func (r *Repos) GetUserByID(id int64) (*model.User, error) {
	var u model.User
	if err := r.DB.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repos) UpdateUserStatus(id int64, status int16) error {
	return r.DB.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}
