package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/codevault-llc/humblebrag-api/config"
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

// Get user by ID
func GetUserById(id uint) (models.UserModel, error) {
	var user models.UserModel

	if err := database.DB.Where("id = ?", id).
		Preload("Subscriptions", func(db *gorm.DB) *gorm.DB {
			return db.Order("updated_at DESC").Limit(3)
		}).
		Preload("Scans", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(5)
		}).
		First(&user).
		Error; err != nil {
		return user, err
	}

	return user, nil
}

// Get user by email
func GetUserByEmail(email string) (models.UserModel, error) {
	var user models.UserModel

	if err := database.DB.Where("email = ?", email).
		First(&user).
		Error; err != nil {
		return user, err
	}

	return user, nil
}

// Get user by Discord ID
func GetUserByDiscordId(discordId string) (models.UserModel, error) {
	var user models.UserModel

	if err := database.DB.Where("discord_id = ?", discordId).
		First(&user).
		Error; err != nil {
		return user, err
	}

	return user, nil
}

// Create a new user
func CreateUser(user models.UserModel) (models.UserModel, error) {
	if err := database.DB.Create(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

// Get user by token
func IsValidUserToken(token string) (models.UserTokenModel, error) {
	var userToken models.UserTokenModel

	if err := database.DB.Where("token = ?", token).
		First(&userToken).
		Error; err != nil {
		return userToken, err
	}

	return userToken, nil
}

// Get user by token
func UpdateUser(user models.UserModel) (models.UserModel, error) {
	if err := database.DB.Save(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

// Get user by token
func FindOrCreateUserFromDiscord(discordUser utils.DiscordUser, token *oauth2.Token) (models.UserModel, error) {
	user, err := GetUserByDiscordId(fmt.Sprint(discordUser.Id))
	if err != nil {
		return models.UserModel{}, err
	}

	if user.ID != 0 {
		return user, nil
	}

	userModel := models.UserModel{
		DiscordId:        fmt.Sprint(discordUser.Id),
		Username:         discordUser.Username,
		Email:            discordUser.Email,
		Avatar:           discordUser.Avatar,
		AccessToken:      token.AccessToken,
		Provider:         "discord",
		StripeCustomerID: "",
		History:          []models.HistoryModel{},
		Subscriptions:    []models.SubscriptionModel{},
	}

	user, err = CreateUser(userModel)
	if err != nil {
		return models.UserModel{}, err
	}

	return user, nil
}

// Get user by token
func FetchDiscordUserInfo(token oauth2.Token) (*utils.DiscordUser, error) {
	res, err := config.DiscordConfig.Client(context.Background(), &token).Get("https://discord.com/api/users/@me")
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user info: %v", err)
	}
	defer res.Body.Close()

	var discordUser utils.DiscordUser
	if err := json.NewDecoder(res.Body).Decode(&discordUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %v", err)
	}

	return &discordUser, nil
}

// Get user by token
func SaveUserToken(userToken string, userID uint) error {
	token := models.UserTokenModel{
		Token:  userToken,
		UserID: userID,
	}

	if err := database.DB.Create(&token).Error; err != nil {
		return err
	}

	return nil
}

func RemoveUserToken(token string) error {
	if err := database.DB.Where("token = ?", token).Delete(&models.UserTokenModel{}).Error; err != nil {
		return err
	}

	return nil
}

const (
	MaxScansPerDay     = 5
	FreeScansPerDay    = 1
	SubscriptionActive = "active"
	ScanStatusComplete = "complete"
)

func CanPerformScan(subscriptions []models.SubscriptionModel, scans []models.ScanModel) bool {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Count completed scans for today
	completedScansToday := 0
	for _, scan := range scans {
		if scan.CreatedAt.After(startOfDay) && scan.Status == ScanStatusComplete {
			completedScansToday++
		}
	}

	fmt.Println("Completed scans today:", completedScansToday)

	// Check if user has an active subscription
	hasActiveSubscription := false
	if len(subscriptions) > 0 {
		// Assuming the first subscription is the latest one
		latestSubscription := subscriptions[0]
		hasActiveSubscription = latestSubscription.Status == SubscriptionActive
	}

	if hasActiveSubscription {
		return completedScansToday < MaxScansPerDay
	}

	// For non-subscribed users
	return completedScansToday < FreeScansPerDay
}
