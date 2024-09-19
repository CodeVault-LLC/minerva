package models

import "gorm.io/gorm"

type NotificationType string

const (
	NotificationGeneral      NotificationType = "general"
	NotificationSecurity     NotificationType = "security"
	NotificationSubscription NotificationType = "subscription"
)

type NotificationModel struct {
	gorm.Model
	UserID         uint             `gorm:"not null"`      // Reference to the user
	Type           NotificationType `gorm:"not null"`      // Type of notification (e.g., security, subscription)
	Message        string           `gorm:"not null"`      // Notification message
	IsRead         bool             `gorm:"default:false"` // Track if the notification has been read
	AdditionalData string           `gorm:"type:text"`     // Optional additional data (e.g., a URL link)
}

type NotificationResponse struct {
	ID        uint             `json:"id"`
	UserID    uint             `json:"user_id"`
	Type      NotificationType `json:"type"`
	Message   string           `json:"message"`
	IsRead    bool             `json:"is_read"`
	CreatedAt string           `json:"created_at"`
}

func ConvertNotifications(notifications []NotificationModel) []NotificationResponse {
	var notificationResponses []NotificationResponse

	for _, notification := range notifications {
		notificationResponses = append(notificationResponses, NotificationResponse{
			ID:        notification.ID,
			UserID:    notification.UserID,
			Type:      notification.Type,
			Message:   notification.Message,
			IsRead:    notification.IsRead,
			CreatedAt: notification.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return notificationResponses
}

func ConvertNotification(notification NotificationModel) NotificationResponse {
	return NotificationResponse{
		ID:        notification.ID,
		UserID:    notification.UserID,
		Type:      notification.Type,
		Message:   notification.Message,
		IsRead:    notification.IsRead,
		CreatedAt: notification.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
