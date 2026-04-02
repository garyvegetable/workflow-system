package repository

import (
	"workflow-system/internal/domain/notification"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notif *notification.Notification) error {
	return r.db.Create(notif).Error
}

func (r *NotificationRepository) ListByUser(userID int64) ([]notification.Notification, error) {
	var notifications []notification.Notification
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) MarkAsRead(id int64) error {
	return r.db.Model(&notification.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *NotificationRepository) MarkAllAsRead(userID int64) error {
	return r.db.Model(&notification.Notification{}).Where("user_id = ?", userID).Update("is_read", true).Error
}

func (r *NotificationRepository) GetUnreadCount(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&notification.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}
