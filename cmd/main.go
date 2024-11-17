package main

import (
	"fmt"

	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/config"
	"github.com/nuhmanudheent/hosp-connect-notification-service/internal/di"
)

func main() {
	config.LoadEnv()
	fmt.Println("Notification service is running and waiting for Kafka events...")
	
	di.KafkaSetup()
	select {}
}
