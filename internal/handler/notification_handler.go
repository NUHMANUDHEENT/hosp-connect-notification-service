package handler

import "github.com/nuhmanudheent/hosp-connect-notification-service/internal/service"

type notificationHandler struct {
	service service.NotificationService
}

func NewNotificationHandler(service service.NotificationService) *notificationHandler {
	return &notificationHandler{
		service: service,
	}
}
func (n *notificationHandler) CunsumeHandler(topic string) error {
	return n.service.SubscribeAndConsume(topic)
}
