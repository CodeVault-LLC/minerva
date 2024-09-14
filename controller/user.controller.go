package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/utils"
	"golang.org/x/oauth2"
)

func GetUserById(id uint) (models.User, error) {
	var user models.User

	if err := constants.DB.Where("id = ?", id).
		First(&user).
		Error; err != nil {
		return user, err
	}

	return user, nil
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User

	if err := constants.DB.Where("email = ?", email).
		First(&user).
		Error; err != nil {
		return user, err
	}

	return user, nil
}

func GetUserByDiscordId(discordId string) (models.User, error) {
	var user models.User

	if err := constants.DB.Where("discord_id = ?", discordId).
		First(&user).
		Error; err != nil {
		return user, err
	}

	return user, nil
}

func CreateUser(user models.User) (models.User, error) {
	if err := constants.DB.Create(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func IsValidUserToken(token string) (models.UserToken, error) {
	var userToken models.UserToken

	if err := constants.DB.Where("token = ?", token).
		First(&userToken).
		Error; err != nil {
		return userToken, err
	}

	return userToken, nil
}

func UpdateUser(user models.User) (models.User, error) {
	if err := constants.DB.Save(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func FindOrCreateUserFromDiscord(discordUser utils.DiscordUser, token *oauth2.Token) (models.User, error) {
	user := models.User{}
	constants.DB.Where("discord_id = ?", discordUser.Id).First(&user)

	if user.ID == 0 {
		user = models.User{
			DiscordId:        fmt.Sprint(discordUser.Id),
			Username:         discordUser.Username,
			Email:            discordUser.Email,
			Avatar:           discordUser.Avatar,
			AccessToken:      token.AccessToken,
			Provider:         "discord",
			StripeCustomerID: "",
			History:          []models.History{},
			Subscriptions:    []models.Subscription{},
			Scans:            []models.Scan{},
		}

		user, err := CreateUser(user)
		if err != nil {
			return user, err
		}
	}

	return user, nil
}

func FetchDiscordUserInfo(token oauth2.Token) (*utils.DiscordUser, error) {
	res, err := constants.DiscordConfig.Client(context.Background(), &token).Get("https://discord.com/api/users/@me")
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

func SaveUserToken(userToken string, userID uint) error {
	token := models.UserToken{
		Token:  userToken,
		UserID: userID,
	}

	if err := constants.DB.Create(&token).Error; err != nil {
		return err
	}

	return nil
}

func RemoveUserToken(token string) error {
	if err := constants.DB.Where("token = ?", token).Delete(&models.UserToken{}).Error; err != nil {
		return err
	}

	return nil
}
