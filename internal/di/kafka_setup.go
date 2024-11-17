package di

import (
	"log"
	"os"
	"time"

	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/config"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/handler"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/repository"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/service"
	"github.com/segmentio/kafka-go"
)

func NewKafkaConsumer(broker, groupID, topic string) (*kafka.Reader, error) {

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{broker},
		GroupID:     groupID,
		Topic:       topic,
		StartOffset: kafka.FirstOffset, 
		MinBytes:    10e3,           
		MaxBytes:    10e6,            
		MaxWait:     1 * time.Second,  
	})

	return reader, nil
}

func KafkaSetup() {

	db := config.InitDatabase()
	repo := repository.NewNotificationRepo(db)

	paymentConsumer, err := NewKafkaConsumer(os.Getenv("KAFKA_BROKER"), "payment_service_group", "payment_topic")
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer for payment topic: %v", err)
	}

	appointmentConsumer, err := NewKafkaConsumer(os.Getenv("KAFKA_BROKER"), "appointment_service_group", "appointment_topic")
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer for appointment topic: %v", err)
	}

	appointmentAlertConsumer, err := NewKafkaConsumer(os.Getenv("KAFKA_BROKER"), "appointment_alert_service_group", "alert_topic")
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer for alert topic: %v", err)
	}


	srv := service.NewNotificationService(repo, paymentConsumer, appointmentConsumer, appointmentAlertConsumer)
	handl := handler.NewNotificationHandler(srv)

	go func() {
		err := handl.PaymentHandler()
		if err != nil {
			log.Fatalf("Error in payment consumer: %v", err)
		}
	}()
	
	go func() {
		err := handl.AppointmentHandler()
		if err != nil {
			log.Fatalf("Error in appointment consumer: %v", err)
		}
	}()

	go func() {
		err := handl.AppointmentAlertHandler()
		if err != nil {
			log.Fatalf("Error in alert_topic consumer: %v", err)
		}
	}()
}
