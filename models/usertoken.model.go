package models

import "gorm.io/gorm"

type UserTokenModel struct {
	gorm.Model

	Token string

	UserID uint
	User   UserModel
}

type UserTokenResponse struct {
	ID uint `json:"id"`

	Token string `json:"token"`

	UserID uint `json:"userId"`

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
