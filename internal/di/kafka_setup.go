package di

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewKafkaConsumer(broker string) (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          "notification_service_group", 
		"auto.offset.reset": "earliest",                   
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	return consumer, nil
}
