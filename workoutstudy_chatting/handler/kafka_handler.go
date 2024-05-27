package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"workoutstudy_chatting/model"
	"workoutstudy_chatting/service"

	"github.com/segmentio/kafka-go"
)

// FitMateHandler는 'fit-mate' 토픽의 메시지를 컨슘
func FitMateHandler(msg kafka.Message, fitMateService service.FitMateService, fitGroupService service.FitGroupServiceInterface) {
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

	var apiResponse model.GetFitMatesApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("Error unmarshalling API response: %v\n", err)
		return
	}

	// API 응답을 서비스 레이어로 넘겨 실제 비즈니스 로직 수행
	if err := fitMateService.HandleFitMateEvent(apiResponse); err != nil {
		log.Printf("Error handling fit mate event: %v\n", err)
	}
}

// FitGroupHandler는 'fit-group' 토픽의 메시지를 컨슘
func FitGroupHandler(msg kafka.Message, fgService service.FitGroupServiceInterface) {
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

	var apiResponse model.GetFitGroupDetailApiResponse
	if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
		log.Printf("Error decoding API response: %v\n", err)
		return
	}

	if err := fgService.HandleFitGroupEvent(apiResponse); err != nil {
		log.Printf("Error handling fit group event: %v\n", err)
	}
}

// HandleMessage는 토픽에 따라 적절한 핸들러를 호출합니다.
func HandleMessage(msg kafka.Message, fitGroupService *service.FitGroupService, fitMateService service.FitMateService) {
	switch msg.Topic {
	case "fit-mate":
		FitMateHandler(msg, fitMateService, fitGroupService)
	case "fit-group":
		FitGroupHandler(msg, fitGroupService)
	default:
		fmt.Printf("No handler for topic %s\n", msg.Topic)
	}
}
