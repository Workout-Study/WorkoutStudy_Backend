package config

import (
	"context"
	"fmt"
	"log"

	"workoutstudy_chatting/handler"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	Readers map[string]*kafka.Reader // 토픽 별로 Reader 저장
}

// KafkaConsumer 생성자
func NewKafkaConsumer(bootstrapServers string, groupID string, topics []string) *KafkaConsumer {
	// Kafka 브로커에 명시적으로 연결 시도
	if err := checkKafkaConnection(bootstrapServers); err != nil {
		log.Fatalf("Failed to connect to Kafka broker at %s: %v", bootstrapServers, err)
	}
	log.Printf("Successfully connected to Kafka broker at %s", bootstrapServers)

	readers := make(map[string]*kafka.Reader)
	for _, topic := range topics {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{bootstrapServers},
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
			// CommitInterval: time.Second,
		})
		readers[topic] = reader
		log.Printf("Kafka Reader created for topic: %s", topic) // 토픽별 Kafka Reader 생성 로그 추가
	}
	return &KafkaConsumer{
		Readers: readers,
	}
}

// Kafka 브로커 연결 확인 함수
func checkKafkaConnection(bootstrapServers string) error {
	// Kafka 브로커에 연결 시도
	conn, err := kafka.Dial("tcp", bootstrapServers)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 연결 성공 시 파티션 정보를 가져와 확인 (명시적 확인)
	partitions, err := conn.ReadPartitions()
	if err != nil {
		return fmt.Errorf("failed to read partitions: %w", err)
	}
	if len(partitions) == 0 {
		return fmt.Errorf("no partitions found")
	}

	return nil
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
