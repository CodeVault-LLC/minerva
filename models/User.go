package models

import "gorm.io/gorm"

type RoleEnum string

const (
	RoleAdmin RoleEnum = "admin"
	RoleUser  RoleEnum = "user"
)

type User struct {
	gorm.Model

	DiscordId string `gorm:"unique;not null"`

	Username string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Avatar   string `gorm:"not null"`

	Role RoleEnum `gorm:"not null" default:"user"`

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

	Role RoleEnum `json:"role"`

	Subscription SubscriptionResponse `json:"subscription"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserMinimalResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
