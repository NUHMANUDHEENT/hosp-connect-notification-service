package handler

import (
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/service"
)

type notificationHandler struct {
	service service.NotificationService
}

func NewNotificationHandler(service service.NotificationService) *notificationHandler {
	return &notificationHandler{
		service: service,
	}
}
func (n *notificationHandler) PaymentHandler() error {
	return n.service.SubscribeAndConsumePaymentEvents()
}

func (n *notificationHandler) AppointmentHandler() error {
	return n.service.SubscribeAndConsumeVideoAppointments()
}

func (n *notificationHandler) AppointmentAlertHandler() error {
	return n.service.SubscribeAndConsumeAppointmentAlerts()
}
