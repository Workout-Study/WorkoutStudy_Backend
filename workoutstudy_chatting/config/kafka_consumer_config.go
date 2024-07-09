package config

import (
	"context"
	"log"
	"time"
	"workoutstudy_chatting/handler"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	Readers map[string]*kafka.Reader // 토픽 별로 Reader 저장
}

// KafkaConsumer 생성자
func NewKafkaConsumer(bootstrapServers []string, groupID string, topics []string) *KafkaConsumer {
	readers := make(map[string]*kafka.Reader)
	for _, topic := range topics {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:        bootstrapServers,
			GroupID:        groupID,
			Topic:          topic,
			MinBytes:       10e3, // 10KB
			MaxBytes:       10e6, // 10MB
			CommitInterval: 0,    // 자동 커밋 비활성화
		})
		readers[topic] = reader
		log.Printf("Kafka Reader created for topic: %s", topic)
	}
	return &KafkaConsumer{
		Readers: readers,
	}
}

// 메시지 Consume 메서드
func (kc *KafkaConsumer) Consume(ctx context.Context, msgChan chan handler.MessageEvent) {
	for topic, reader := range kc.Readers {
		go func(topic string, r *kafka.Reader) {
			log.Printf("Starting Kafka Consumer for topic: %s", topic)
			for { 
				m, err := r.FetchMessage(ctx)
				if err != nil {
					log.Printf("Error fetching message from topic %s: %v\n", topic, err)
					if err == context.Canceled {
						return
					}
					time.Sleep(time.Second) // 재시도 전에 잠시 대기
					continue
				}
				log.Printf("Message received from topic %s: %s\n", topic, string(m.Value))
				msgChan <- handler.MessageEvent{Message: m, Topic: topic}

				// 명시적으로 커밋
				if err := r.CommitMessages(ctx, m); err != nil {
					log.Printf("Failed to commit message for topic %s: %v\n", topic, err)
				}
			}
		}(topic, reader)
	}
}
