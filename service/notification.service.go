package service

import (
	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/models"
)

// CreateNotification creates a new notification in the database
func CreateNotification(notification models.Notification) (*models.Notification, error) {
	tx := constants.DB.Begin()
	if err := tx.Create(&notification).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &notification, tx.Commit().Error
}

// GetNotificationsByUserID retrieves all notifications for a user
func GetNotificationsByUserID(userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	if err := constants.DB.Where("user_id = ?", userID).Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}

// GetUnreadNotificationsByUserID retrieves all unread notifications for a user
func GetUnreadNotificationsByUserID(userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	if err := constants.DB.Where("user_id = ? AND is_read = ?", userID, false).Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}

// MarkNotificationAsRead updates the read status of a notification
func MarkNotificationAsRead(notificationID uint) error {
	tx := constants.DB.Begin()
	if err := tx.Model(&models.Notification{}).Where("id = ?", notificationID).Update("is_read", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
