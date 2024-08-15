package controller

import (
	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/models"
)

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
