package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/domain"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/repository"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/utils"
)

type PaymentEvent struct {
	PaymentID string  `json:"payment_id"`
	PatientID string  `json:"patient_id"`
	Email     string  `json:"email"`
	Amount    float64 `json:"amount"`
	Date      string  `json:"date"`
}

type NotificationService interface {
	PaymentSubscribeAndConsume(topic string) error
	VideoAppointmentSubcribeAndCunsume(topic string) error
	AppointmentAlertSubscribeAndConsume(topic string) error
}
type notificationService struct {
	Paymentconsumer          *kafka.Consumer
	VideoAppointmentConsumer *kafka.Consumer
	AppointmentAlertConsumer *kafka.Consumer
	repo                     repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository, paymentcons, videoappocons, appalertcons *kafka.Consumer) NotificationService {
	return &notificationService{
		repo:                     repo,
		Paymentconsumer:          paymentcons,
		VideoAppointmentConsumer: videoappocons,
		AppointmentAlertConsumer: appalertcons,
	}
}

// SubscribeAndConsume listens to the topic and processes messages
func (kc *notificationService) PaymentSubscribeAndConsume(topic string) error {
	fmt.Println("topic:", topic)
	err := kc.Paymentconsumer.Subscribe(topic, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	for {
		msg, err := kc.Paymentconsumer.ReadMessage(-1)
		if err != nil {
			fmt.Println("error", err)
			continue
		}
		var paymentEvent PaymentEvent
		// Deserialize the message
		err = json.Unmarshal(msg.Value, &paymentEvent)
		if err != nil {
			fmt.Printf("Failed to unmarshal message: %v\n", err)
			continue
		}

		fmt.Printf("Received payment event: %+v\n", paymentEvent)

		// Send notification (e.g., send email)
		notificationData, err := SendNotificationSetup(paymentEvent.Email, paymentEvent)
		if err != nil {
			fmt.Printf("Failed to send notification: %v\n", err)
		}
		err = kc.repo.NotificationStore(notificationData)
		if err != nil {
			fmt.Printf("Failed to store notification: %v\n", err)

		}
		time.Sleep(time.Second)

	}
}

// SendNotification sends an email notification
func SendNotificationSetup(email string, event PaymentEvent) (domain.Notification, error) {
	subject := "Payment Confirmation"
	body := fmt.Sprintf("Dear Patient,\n\n"+
		"Your payment has been successfully processed.\n\n"+
		"Details:\n"+
		"Payment ID: %s\n"+
		"Patient ID: %s\n"+
		"Amount: %.2f\n"+
		"Date: %s\n\n"+
		"Thank you for your payment!",
		event.PaymentID, event.PatientID, event.Amount, event.Date)

	// Call the function to send the email (this would be your email service)
	err := utils.SendNotificationToEmail(email, subject, body)
	if err != nil {
		return domain.Notification{}, fmt.Errorf("failed to send email: %w", err)
	}
	var notificationstore = domain.Notification{
		Message: body,
		UserId:  event.PatientID,
		// Email:   email,
		Status: "send",
	}

	fmt.Println("Payment completion notification sent successfully to:", email)
	return notificationstore, nil
}
func (kc *notificationService) VideoAppointmentSubcribeAndCunsume(topic string) error {
	fmt.Println("topic:", topic)
	err := kc.VideoAppointmentConsumer.Subscribe(topic, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	for {
		msg, err := kc.VideoAppointmentConsumer.ReadMessage(-1)
		fmt.Println("hiiiii")
		if err == nil {
			var appointmentEvent domain.AppointmentEvent
			err = json.Unmarshal(msg.Value, &appointmentEvent)
			if err != nil {
				fmt.Printf("Failed to unmarshal message: %v\n", err)
				continue
			}

			fmt.Printf("Received appointment event: %+v\n", appointmentEvent)

			// Send notification (e.g., send email)
			notificationData, err := SendVideoAppointmentNotification(appointmentEvent.Email, appointmentEvent)
			if err != nil {
				fmt.Printf("Failed to send notification: %v\n", err)
			}
			err = kc.repo.NotificationStore(notificationData)
			if err != nil {
				fmt.Printf("Failed to store notification: %v\n", err)
			}
		} else {
			fmt.Printf("Consumer error: %v\n", err)
		}
		time.Sleep(time.Second)

	}
}

func SendVideoAppointmentNotification(email string, event domain.AppointmentEvent) (domain.Notification, error) {
	fmt.Println(event)
	subject := "Video Appointment Meet Link"
	body := fmt.Sprintf("Dear Patient,\n\n"+
		"This is a reminder for your video appointment right now.\n\n"+
		"Details:\n"+
		"Appointment ID: %d\n"+
		"Doctor: Dr. %s\n"+
		"Appointment Date And Time: %s\n"+
		"Please join the video call using the following link: %s\n\n"+
		"Thank you!",
		event.AppointmentId, event.DoctorId, event.AppointmentDate, event.VideoURL)

	err := utils.SendNotificationToEmail(email, subject, body)
	if err != nil {
		return domain.Notification{}, fmt.Errorf("failed to send email: %w", err)
	}

	notificationstore := domain.Notification{
		Message: body,
		UserId:  event.Email,
		Status:  "send",
	}

	fmt.Println("Video appointment notification sent successfully to:", email)
	return notificationstore, nil
}
func (kc *notificationService) AppointmentAlertSubscribeAndConsume(topic string) error {
	fmt.Println("topic:", topic)
	err := kc.AppointmentAlertConsumer.Subscribe(topic, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	for {
		msg, err := kc.AppointmentAlertConsumer.ReadMessage(-1)
		if err == nil {
			var appointmentEvent domain.AppointmentEvent
			err = json.Unmarshal(msg.Value, &appointmentEvent)
			if err != nil {
				fmt.Printf("Failed to unmarshal message: %v\n", err)
				continue
			}

			fmt.Printf("Received appointment alert event: %+v\n", appointmentEvent)

			// Send notification (e.g., send email)
			notificationData, err := SendAppointmentAlertNotification(appointmentEvent.Email, appointmentEvent)
			if err != nil {
				fmt.Printf("Failed to send alert notification: %v\n", err)
			}
			err = kc.repo.NotificationStore(notificationData)
			if err != nil {
				fmt.Printf("Failed to store alert notification: %v\n", err)
			}
		} else {
			fmt.Printf("Consumer error: %v\n", err)
		}
		time.Sleep(time.Second)
	}
}


func SendAppointmentAlertNotification(email string, event domain.AppointmentEvent) (domain.Notification, error) {
	fmt.Println(event)
	subject := "Appoinmtnt Alert"
	body := fmt.Sprintf("Dear Patient,\n\n"+
		"This is a reminder for your appointment on today.\n\n"+
		"Details:\n"+
		"Appointment ID: %d\n"+
		"Doctor: Dr. %s\n"+
		"Appointment Date And Time: %s\n"+
		"Please be available the sheduled time: %s\n\n"+
		"Thank you!",
		event.AppointmentId, event.DoctorId, event.AppointmentDate, event.VideoURL)

	err := utils.SendNotificationToEmail(email, subject, body)
	if err != nil {
		return domain.Notification{}, fmt.Errorf("failed to send email: %w", err)
	}

	notificationstore := domain.Notification{
		Message: body,
		UserId:  event.Email,
		Status:  "send",
	}

	fmt.Println("alert appointment notification sent successfully to:", email)
	return notificationstore, nil
}
