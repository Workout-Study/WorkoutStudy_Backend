package config

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// KafkaConsumer 구조체 정의
type KafkaConsumer struct {
	Consumer *kafka.Consumer // 소문자 'consumer'에서 대문자 'Consumer'로 변경하여 공개 필드로 만듦
}

// NewKafkaConsumer 함수는 KafkaConsumer 인스턴스를 초기화합니다.
func NewKafkaConsumer(bootstrapServers string) *KafkaConsumer {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":        bootstrapServers,
		"group.id":                 "chatting-server-consumer",
		"auto.offset.reset":        "latest",
		"enable.auto.commit":       false,
		"isolation.level":          "read_committed",
		"allow.auto.create.topics": false,
	}

	consumer, err := kafka.NewConsumer(configMap)
	if err != nil {
		panic(fmt.Sprintf("Failed to create consumer: %s", err))
	}

	return &KafkaConsumer{
		Consumer: consumer,
	}
}

func (kc *KafkaConsumer) Consume(topics []string) {
	err := kc.Consumer.SubscribeTopics(topics, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to subscribe to topics: %s", err))
	}

	fmt.Println("Kafka consumer started. Waiting for messages...")
	for {
		msg, err := kc.Consumer.ReadMessage(-1)
		if err != nil {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			continue
		}
		fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
	}
}
