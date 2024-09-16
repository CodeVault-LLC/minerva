package constants

import (
	"fmt"

	"github.com/codevault-llc/humblebrag-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(DATABASE_URL string) {
	db, err := gorm.Open(postgres.Open(DATABASE_URL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	registerGlobalCallbacks(db)

	db.AutoMigrate(&models.User{}, &models.Subscription{}, &models.Scan{}, &models.Detail{}, &models.Finding{}, &models.Certificate{}, &models.History{}, &models.UserToken{}, &models.Content{}, &models.List{})
	DB = db

	fmt.Println("Connected to PostgreSQL")
}

func handleRecordNotFound(db *gorm.DB) {
	if db.Error == gorm.ErrRecordNotFound {
		db.Error = nil
	}
}

func registerGlobalCallbacks(db *gorm.DB) {
	db.Callback().Query().After("gorm:query").Register("app:handle_record_not_found", handleRecordNotFound)
}
