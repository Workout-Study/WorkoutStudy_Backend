package config

import (
	"context"
	"log"

	"workoutstudy_chatting/handler"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	Readers map[string]*kafka.Reader // 토픽 별로 Reader 저장
}

// KafkaConsumer 생성자
func NewKafkaConsumer(bootstrapServers string, groupID string, topics []string) *KafkaConsumer {
	readers := make(map[string]*kafka.Reader)
	for _, topic := range topics {
		readers[topic] = kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{bootstrapServers},
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		})
		log.Printf("Kafka Reader created for topic: %s", topic) // 토픽별 Kafka Reader 생성 로그 추가
	}
	return &KafkaConsumer{
		Readers: readers,
	}
}

// 메시지 Consume 메서드
func (kc *KafkaConsumer) Consume(ctx context.Context, msgChan chan handler.MessageEvent) {
	for topic, reader := range kc.Readers {
		go func(topic string, r *kafka.Reader) {
			log.Printf("Starting Kafka Consumer for topic: %s", topic) // 컨슈머 시작 로그 추가
			for {
				m, err := r.ReadMessage(ctx)
				if err != nil {
					log.Printf("Error reading message from topic %s: %v\n", topic, err)
					break
				}
				log.Printf("Message received from topic %s: %s\n", topic, string(m.Value)) // 디버깅 로그 추가
				msgChan <- handler.MessageEvent{Message: m, Service: nil}
			}
		}(topic, reader)
	}
}
