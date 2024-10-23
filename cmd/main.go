package main

import (
	"fmt"
	"log"

	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/config"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/di"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/handler"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/repository"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/service"
)

func main() {
	di.LoadEnv()

	db := config.InitDatabase()
	repo := repository.NewNotificationRepo(db)
	// Initialize Kafka consumer
	kafkaConsumer, err := di.NewKafkaConsumer("localhost:9092")
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer kafkaConsumer.Close()
	srv := service.NewNotificationService(repo, kafkaConsumer)
	handl := handler.NewNotificationHandler(srv)

	// go func() {
	// 	err := handl.PaymentHandler("payment_topic")
	// 	if err != nil {
	// 		log.Fatalf("Error in Kafka consumer: %v", err)
	// 	}
	// }()
	go func() {
		err := handl.AppointmentHandler("appointment_topic")
		if err != nil {
			log.Fatalf("Error in Kafka consumer: %v", err)
		}
	}()

	fmt.Println("Notification service is running and waiting for Kafka events...")

	// Keep the main function alive
	select {}
}
