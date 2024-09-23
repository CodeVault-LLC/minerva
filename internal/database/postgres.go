package database

import (
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"

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

	err = db.AutoMigrate(&models.UserModel{}, &models.NotificationModel{}, &models.SubscriptionModel{}, &models.ScanModel{}, &models.NetworkModel{}, &models.WhoisModel{}, &models.FindingModel{}, &models.CertificateModel{}, &models.CertificateResultModel{}, &models.HistoryModel{}, &models.UserTokenModel{}, &models.ContentModel{}, &models.ListModel{})
	if err != nil {
		logger.Log.Error("Failed to auto migrate models: %v", err)
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
		logger.Log.Error("Failed to register global callbacks: %v", err)
	}
}
