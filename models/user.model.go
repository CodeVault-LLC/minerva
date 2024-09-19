package models

import "gorm.io/gorm"

type RoleEnum string

const (
	RoleAdmin RoleEnum = "admin"
	RoleUser  RoleEnum = "user"
)

type TwoFAMethod string

const (
	MethodApp   TwoFAMethod = "app"   // e.g., Google Authenticator
	MethodEmail TwoFAMethod = "email" // e.g., email-based OTP
)

type UserModel struct {
	gorm.Model

	DiscordId string `gorm:"unique;not null"`

	Username string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Avatar   string `gorm:"not null"`

	Role RoleEnum `gorm:"not null" default:"user"`

	AccessToken string `gorm:"not null"`

	Provider         string `gorm:"not null"`
	StripeCustomerID string `gorm:"not null"`

	// 2FA-related fields
	Is2FAEnabled bool        `gorm:"default:false"` // Track if 2FA is enabled
	TwoFASecret  string      `gorm:"not null"`      // Secret key for 2FA
	TwoFAMethod  TwoFAMethod `gorm:"default:'app'"` // 2FA method, app-based or email-based

	History       []HistoryModel      `gorm:"foreignKey:UserID"`
	Subscriptions []SubscriptionModel `gorm:"foreignKey:UserID"`
	Scans         []ScanModel         `gorm:"foreignKey:UserID"`
	Notifications []NotificationModel `gorm:"foreignKey:UserID"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	AvatarURL string `json:"avatar_url"`

	Role RoleEnum `json:"role"`

	Subscription  SubscriptionResponse   `json:"subscription"`
	Notifications []NotificationResponse `json:"notifications"`

	Is2FAEnabled bool `json:"is_2fa_enabled"` // Reflect 2FA status

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserMinimalResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func ConvertUser(user UserModel) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Avatar:    user.Avatar,
		AvatarURL: "https://cdn.discordapp.com/avatars/" + user.DiscordId + "/" + user.Avatar + ".png",
		CreatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ConvertUserMinimal(user UserModel) UserMinimalResponse {
	return UserMinimalResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}