package repository

import (
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"gorm.io/gorm"
)

type RedirectRepo struct {
	db *gorm.DB
}

func NewRedirectRepo(db *gorm.DB) *RedirectRepo {
	return &RedirectRepo{db: db}
}

var RedirectRepository *RedirectRepo

func (r *RedirectRepo) Create(redirect entities.RedirectModel) (entities.RedirectModel, error) {
	tx := r.db.Begin()
	if err := tx.Create(&redirect).Error; err != nil {
		tx.Rollback()
		return entities.RedirectModel{}, err
	}
	tx.Commit()
	return redirect, nil
}
