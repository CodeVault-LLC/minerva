package repository

import (
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"gorm.io/gorm"
)

type DnsRepo struct {
	db *gorm.DB
}

func NewDnsRepository(db *gorm.DB) *DnsRepo {
	return &DnsRepo{
		db: db,
	}
}

var DnsRepository *DnsRepo

func (repository *DnsRepo) SaveDnsResult(dns entities.DNSModel) error {
	tx := repository.db.Begin()
	if err := tx.Create(&dns).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
