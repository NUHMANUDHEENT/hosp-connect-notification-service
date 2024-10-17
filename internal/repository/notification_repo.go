package repository

import (
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/domain"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	NotificationStore(noti domain.Notification) error
}
type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepo(db *gorm.DB) NotificationRepository {
	return &notificationRepository{
		db: db,
	}
}
func (r *notificationRepository) NotificationStore(noti domain.Notification) error {
	return r.db.Create(&noti).Error
}
