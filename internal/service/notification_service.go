package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/domain"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/repository"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/utils"
	"github.com/segmentio/kafka-go"
)

// In service/notificationService.go
type NotificationService interface {
	SubscribeAndConsumePaymentEvents() error
	SubscribeAndConsumeVideoAppointments() error
	SubscribeAndConsumeAppointmentAlerts() error
}
type notificationService struct {
	paymentConsumer          *kafka.Reader
	videoAppointmentConsumer *kafka.Reader
	appointmentAlertConsumer *kafka.Reader
	repo                     repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository, paymentCons, videoApptCons, alertCons *kafka.Reader) NotificationService {
	return &notificationService{
		repo:                     repo,
		paymentConsumer:          paymentCons,
		videoAppointmentConsumer: videoApptCons,
		appointmentAlertConsumer: alertCons,
	}
}

func (ns *notificationService) SubscribeAndConsumePaymentEvents() error {
	for {
		msg, err := ns.paymentConsumer.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Consumer error:", err)
			continue
		}

		var paymentEvent domain.PaymentEvent
		if err := json.Unmarshal(msg.Value, &paymentEvent); err != nil {
			fmt.Printf("Failed to unmarshal payment message: %v\n", err)
			continue
		}

		notification, err := createPaymentNotification(paymentEvent)
		if err != nil {
			fmt.Printf("Failed to create payment notification: %v\n", err)
			continue
		}

		if err := ns.repo.NotificationStore(notification); err != nil {
			fmt.Printf("Failed to store payment notification: %v\n", err)
		}

		// Commit the message after processing
		if err := ns.paymentConsumer.CommitMessages(context.Background(), msg); err != nil {
			fmt.Printf("Failed to commit message: %v\n", err)
		}
		time.Sleep(time.Second)
	}
}

// createPaymentNotification handles sending the email and structuring the notification for storage
func createPaymentNotification(event domain.PaymentEvent) (domain.Notification, error) {
	subject := "Payment Confirmation"
	body := fmt.Sprintf("Dear Patient,\n\nYour payment has been processed.\n\nDetails:\nPayment ID: %s\nAmount: %.2f\nDate: %s\n\nThank you!",
		event.PaymentID, event.Amount, event.Date)

	if err := utils.SendNotificationToEmail(event.Email, subject, body); err != nil {
		return domain.Notification{}, fmt.Errorf("failed to send email: %w", err)
	}

	return domain.Notification{
		Message: body,
		UserId:  event.PatientID,
		Status:  "sent",
	}, nil
}
func (ns *notificationService) SubscribeAndConsumeVideoAppointments() error {
	for {
		msg, err := ns.videoAppointmentConsumer.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Consumer error:", err)
			continue
		}

		var apptEvent domain.AppointmentEvent
		if err := json.Unmarshal(msg.Value, &apptEvent); err != nil {
			fmt.Printf("Failed to unmarshal appointment message: %v\n", err)
			continue
		}

		notification, err := createVideoAppointmentNotification(apptEvent)
		if err != nil {
			fmt.Printf("Failed to create video appointment notification: %v\n", err)
			continue
		}

		if err := ns.repo.NotificationStore(notification); err != nil {
			fmt.Printf("Failed to store video appointment notification: %v\n", err)
		}

		if err := ns.videoAppointmentConsumer.CommitMessages(context.Background(), msg); err != nil {
			fmt.Printf("Failed to commit message: %v\n", err)
		}
		time.Sleep(time.Second)
	}
}

// createVideoAppointmentNotification structures the notification for a video appointment
func createVideoAppointmentNotification(event domain.AppointmentEvent) (domain.Notification, error) {
	subject := "Video Appointment Reminder"
	body := fmt.Sprintf("Dear Patient,\n\nYour video appointment is scheduled.\n\nDetails:\nAppointment ID: %d\nDoctor: Dr. %s\nDate & Time: %s\nLink: %s\n\nThank you!",
		event.AppointmentId, event.DoctorId, event.AppointmentDate, event.VideoURL)

	if err := utils.SendNotificationToEmail(event.Email, subject, body); err != nil {
		return domain.Notification{}, fmt.Errorf("failed to send email: %w", err)
	}

	return domain.Notification{
		Message: body,
		UserId:  event.Email,
		Status:  "sent",
	}, nil
}
func (ns *notificationService) SubscribeAndConsumeAppointmentAlerts() error {
	for {
		msg, err := ns.appointmentAlertConsumer.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Consumer error:", err)
			continue
		}

		var alertEvent domain.AppointmentEvent
		if err := json.Unmarshal(msg.Value, &alertEvent); err != nil {
			fmt.Printf("Failed to unmarshal alert message: %v\n", err)
			continue
		}

		notification, err := createAppointmentAlertNotification(alertEvent)
		if err != nil {
			fmt.Printf("Failed to create appointment alert notification: %v\n", err)
			continue
		}

		if err := ns.repo.NotificationStore(notification); err != nil {
			fmt.Printf("Failed to store appointment alert notification: %v\n", err)
		}

		if err := ns.appointmentAlertConsumer.CommitMessages(context.Background(), msg); err != nil {
			fmt.Printf("Failed to commit message: %v\n", err)
		}
		time.Sleep(time.Second)
	}
}

// createAppointmentAlertNotification structures the alert notification for upcoming appointments
func createAppointmentAlertNotification(event domain.AppointmentEvent) (domain.Notification, error) {
	subject := "Appointment Alert"
	body := fmt.Sprintf("Dear Patient,\n\nYour appointment is scheduled.\n\nDetails:\nAppointment ID: %d\nDoctor: Dr. %s\nDate & Time: %s\nPlease be ready at the scheduled time.\n\nThank you!",
		event.AppointmentId, event.DoctorId, event.AppointmentDate)

	if err := utils.SendNotificationToEmail(event.Email, subject, body); err != nil {
		return domain.Notification{}, fmt.Errorf("failed to send email: %w", err)
	}

	return domain.Notification{
		Message: body,
		UserId:  event.Email,
		Status:  "sent",
	}, nil
}
