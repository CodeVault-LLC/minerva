package database

import (
	contentEntities "github.com/codevault-llc/humblebrag-api/internal/contents/models/entities"
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	networkEntities "github.com/codevault-llc/humblebrag-api/internal/network/models/entities"
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

	err = db.AutoMigrate(&entities.LicenseModel{}, &entities.ScanModel{}, &networkEntities.NetworkModel{},
		&networkEntities.DNSModel{}, &entities.MetadataModel{}, &networkEntities.WhoisModel{}, &contentEntities.FindingModel{},
		&networkEntities.CertificateModel{}, &contentEntities.ContentModel{}, &contentEntities.ContentStorageModel{},
		&contentEntities.ContentTagsModel{}, &contentEntities.ContentAccessLogModel{},
		&entities.FilterModel{}, &entities.RedirectModel{}, &entities.ScreenshotModel{})
	if err != nil {
		logger.Log.Error("Failed to auto migrate entities: %v", zap.Error(err))
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
