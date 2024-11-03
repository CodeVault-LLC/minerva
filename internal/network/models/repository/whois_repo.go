package repository

import (
	"github.com/codevault-llc/humblebrag-api/internal/network/models/entities"
	"gorm.io/gorm"
)

type WhoisRepo struct {
	db *gorm.DB
}

func NewWhoisRepository(db *gorm.DB) *WhoisRepo {
	return &WhoisRepo{
		db: db,
	}
}

var WhoisRepository *WhoisRepo

func (repository *WhoisRepo) SaveWhoisResult(whois entities.WhoisModel) error {
	tx := repository.db.Begin()
	if err := tx.Create(&whois).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
