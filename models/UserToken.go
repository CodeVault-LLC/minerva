package models

import "gorm.io/gorm"

type UserToken struct {
	gorm.Model

	Token string

	UserID uint
	User   User
}

type UserTokenResponse struct {
	ID uint `json:"id"`

	Token string `json:"token"`

	UserID uint `json:"userId"`

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
