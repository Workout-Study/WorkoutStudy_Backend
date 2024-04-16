package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/segmentio/kafka-go"
)

// FitMateHandler는 'fit-mate' 토픽의 메시지를 처리합니다.
func FitMateHandler(msg kafka.Message) {
	value, err := strconv.Atoi(string(msg.Value))
	if err != nil {
		log.Printf("Error converting Kafka message to int: %v\n", err)
		return
	}
	url := fmt.Sprintf("http://fit-group:8080/fit-group-service/groups/%d", value)
	response, err := http.Get(url)
	if err != nil {
		log.Printf("Error sending GET request: %v\n", err)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return
	}
	fmt.Printf("Response from fit-group-service for groups: %s\n", string(body))
}

// FitGroupHandler는 'fit-group' 토픽의 메시지를 처리합니다.
func FitGroupHandler(msg kafka.Message) {
	value, err := strconv.Atoi(string(msg.Value))
	if err != nil {
		log.Printf("Error converting Kafka message to int: %v\n", err)
		return
	}
	url := fmt.Sprintf("http://fit-group:8080/fit-group-service/mates/%d", value)
	response, err := http.Get(url)
	if err != nil {
		log.Printf("Error sending GET request: %v\n", err)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return
	}
	fmt.Printf("Response from fit-group-service for mates: %s\n", string(body))
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
