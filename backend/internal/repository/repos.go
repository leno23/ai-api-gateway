package repository

import (
	"gorm.io/gorm"
)

// Repos holds shared DB access for thin repository methods.
type Repos struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Repos {
	return &Repos{DB: db}
}
