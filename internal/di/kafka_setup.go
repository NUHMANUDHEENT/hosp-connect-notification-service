package di

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/config"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/handler"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/repository"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/service"
)

func NewKafkaConsumer(broker, group string) (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          group,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	return consumer, nil
}
func KafkaSetup() {

	db := config.InitDatabase()

	repo := repository.NewNotificationRepo(db)

	paymentConsumer, err := NewKafkaConsumer("localhost:9092", "payment_service_group")
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer for payment topic: %v", err)
	}

	appointmentConsumer, err := NewKafkaConsumer("localhost:9092", "appointment_service_group")
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer for appointment topic: %v", err)
	}

	appointmentAlertConsumer, err := NewKafkaConsumer("localhost:9092", "appointment_service_group")
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer for appointment alert topic: %v", err)
	}

	defer paymentConsumer.Close()
	defer appointmentConsumer.Close()
	defer appointmentAlertConsumer.Close()

	// Initialize services
	srv := service.NewNotificationService(repo, paymentConsumer, appointmentConsumer, appointmentAlertConsumer)
	handl := handler.NewNotificationHandler(srv)

	// Launch handlers in separate goroutines
	go func() {
		err := handl.PaymentHandler("payment_topic")
		if err != nil {
			log.Fatalf("Error in payment consumer: %v", err)
		}
	}()

	go func() {
		err := handl.AppointmentHandler("appointment_topic")
		if err != nil {
			log.Fatalf("Error in appointment consumer: %v", err)
		}
	}()

	go func() {
		err := handl.AppointmentAlertHandler("alert_topic")
		if err != nil {
			log.Fatalf("Error in alert_topic consumer: %v", err)
		}
	}()
}
