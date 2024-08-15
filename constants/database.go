package constants

import (
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

	db.AutoMigrate(&models.User{}, &models.Subscription{}, &models.Scan{}, &models.Finding{}, &models.History{}, &models.UserToken{}, &models.Content{})
	DB = db
}
