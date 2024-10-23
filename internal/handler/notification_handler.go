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
func (n *notificationHandler) PaymentHandler(topic string) error {
	return n.service.SubscribeAndConsume(topic)
}
func (n *notificationHandler) AppointmentHandler(topic string) error {
	return n.service.VideoAppointmentSubcribeAndCunsume(topic)
}
