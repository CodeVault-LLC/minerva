package utils

import "github.com/codevault-llc/humblebrag-api/models"

func ConvertUser(user models.User) models.UserResponse {
	return models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Avatar:  user.Avatar,
		// use it as the discord avatar url
		AvatarURL: "https://cdn.discordapp.com/avatars/" + user.DiscordId + "/" + user.Avatar + ".png",
		CreatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
