package database

import (
	"github.com/codevault-llc/humblebrag-api/internal/database/models"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"go.uber.org/zap"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitPostgres(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	registerGlobalCallbacks(db)

	err = db.AutoMigrate(&models.LicenseModel{}, &models.ScanModel{}, &models.NetworkModel{},
		&models.DNSModel{}, &models.MetadataModel{}, &models.WhoisModel{}, &models.FindingModel{},
		&models.CertificateModel{}, &models.ContentModel{}, &models.ContentStorageModel{}, &models.ContentTagsModel{}, &models.ContentAccessLogModel{},
		&models.FilterModel{})
	if err != nil {
		logger.Log.Error("Failed to auto migrate models: %v", zap.Error(err))
		return nil, err
	}
	DB = db

	return db, nil
}

func handleRecordNotFound(db *gorm.DB) {
	if db.Error == gorm.ErrRecordNotFound {
		db.Error = nil
	}
}

func registerGlobalCallbacks(db *gorm.DB) {
	err := db.Callback().Query().After("gorm:query").Register("app:handle_record_not_found", handleRecordNotFound)

	if err != nil {
		logger.Log.Error("Failed to register global callbacks: %v", zap.Error(err))
	}
}
