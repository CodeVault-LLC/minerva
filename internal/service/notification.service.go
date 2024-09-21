package service

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/models"
)

// CreateNotification creates a new notification in the database
func CreateNotification(notification models.NotificationModel) (*models.NotificationModel, error) {
	tx := database.DB.Begin()
	if err := tx.Create(&notification).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &notification, tx.Commit().Error
}

// GetNotificationsByUserID retrieves all notifications for a user
func GetNotificationsByUserID(userID uint) ([]models.NotificationModel, error) {
	var notifications []models.NotificationModel
	if err := database.DB.Where("user_id = ?", userID).Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}

// GetUnreadNotificationsByUserID retrieves all unread notifications for a user
func GetUnreadNotificationsByUserID(userID uint) ([]models.NotificationModel, error) {
	var notifications []models.NotificationModel
	if err := database.DB.Where("user_id = ? AND is_read = ?", userID, false).Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}

// MarkNotificationAsRead updates the read status of a notification
func MarkNotificationAsRead(notificationID uint) error {
	tx := database.DB.Begin()
	if err := tx.Model(&models.NotificationModel{}).Where("id = ?", notificationID).Update("is_read", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
