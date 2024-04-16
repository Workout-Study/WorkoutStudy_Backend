package handler

import (
	"fmt"

	"github.com/segmentio/kafka-go"
)

// FitMateHandler는 'fit-mate' 토픽의 메시지를 처리합니다.
func FitMateHandler(msg kafka.Message) {
	// 'fit-mate' 토픽 메시지에 대한 처리 로직
	fmt.Printf("Processing 'fit-mate' message: %s\n", string(msg.Value))
	// 여기에 필요한 비즈니스 로직 추가
}

// FitGroupHandler는 'fit-group' 토픽의 메시지를 처리합니다.
func FitGroupHandler(msg kafka.Message) {
	// 'fit-group' 토픽 메시지에 대한 처리 로직
	fmt.Printf("Processing 'fit-group' message: %s\n", string(msg.Value))
	// 여기에 필요한 비즈니스 로직 추가
}

// HandleMessage는 토픽에 따라 적절한 핸들러를 호출합니다.
func HandleMessage(msg kafka.Message) {
	switch msg.Topic {
	case "fit-mate":
		FitMateHandler(msg)
	case "fit-group":
		FitGroupHandler(msg)
	default:
		fmt.Printf("No handler for topic %s\n", msg.Topic)
	}
}
