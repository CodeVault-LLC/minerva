package database

import (
	"github.com/codevault-llc/humblebrag-api/models"

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

	db.AutoMigrate(&models.UserModel{}, &models.NotificationModel{}, &models.SubscriptionModel{}, &models.ScanModel{}, &models.WhoisModel{}, &models.NetworkModel{}, &models.FindingModel{}, &models.CertificateModel{}, &models.HistoryModel{}, &models.UserTokenModel{}, &models.ContentModel{}, &models.ListModel{})
	DB = db

	return db, nil
}

func handleRecordNotFound(db *gorm.DB) {
	if db.Error == gorm.ErrRecordNotFound {
		db.Error = nil
	}
}

func registerGlobalCallbacks(db *gorm.DB) {
	db.Callback().Query().After("gorm:query").Register("app:handle_record_not_found", handleRecordNotFound)
}
