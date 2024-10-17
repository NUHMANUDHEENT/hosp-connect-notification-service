package service

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/di"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/domain"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/repository"
)

type PaymentEvent struct {
	PaymentID string  `json:"payment_id"`
	PatientID string  `json:"patient_id"`
	Email     string  `json:"email"`
	Amount    float64 `json:"amount"`
	Date      string  `json:"date"`
}

type NotificationService interface {
	SubscribeAndConsume(topic string) error
}
type notificationService struct {
	consumer *kafka.Consumer
	repo     repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository, kafkaConsumer *kafka.Consumer) NotificationService {
	return &notificationService{
		repo:     repo,
		consumer: kafkaConsumer,
	}
}

// SubscribeAndConsume listens to the topic and processes messages
func (kc *notificationService) SubscribeAndConsume(topic string) error {
	fmt.Println(topic)
	err := kc.consumer.Subscribe(topic, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	for {
		msg, err := kc.consumer.ReadMessage(-1) // Blocking call, waits for new messages
		if err == nil {
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
		} else {
			fmt.Printf("Consumer error: %v\n", err)
		}
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
	err := di.SendNotificationToEmail(email, subject, body)
	if err != nil {
		return domain.Notification{}, fmt.Errorf("failed to send email: %w", err)
	}
	var notificationstore = domain.Notification{
		Message: body,
		UserId:  event.PatientID,
		// Email:   email,
		Status: "send",
	}

	fmt.Println("Notification sent successfully to:", email)
	return notificationstore, nil
}
