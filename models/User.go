package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	DiscordId string `gorm:"unique;not null"`

	Username string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Avatar   string `gorm:"not null"`

	AccessToken string `gorm:"not null"`

	Provider         string `gorm:"not null"`
	StripeCustomerID string `gorm:"not null"`

	History       []History      `gorm:"foreignKey:UserID"`
	Subscriptions []Subscription `gorm:"foreignKey:UserID"`
	Scans         []Scan         `gorm:"foreignKey:UserID"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	AvatarURL string `json:"avatar_url"`

	Subscription SubscriptionResponse `json:"subscription"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserMinimalResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
