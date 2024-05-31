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

type MessageEvent struct {
	Message kafka.Message
	Service interface{}
}

func FitMateHandler(c chan MessageEvent, fitMateService service.FitMateService, fitGroupService service.FitGroupServiceInterface) {
	for event := range c {
		msg := event.Message
		value, err := strconv.Atoi(string(msg.Value))
		if err != nil {
			log.Printf("Error converting Kafka message to int: %v\n", err)
			continue
		}

		url := fmt.Sprintf("http://fit-group:8080/fit-group-service/mates/%d", value)
		response, err := http.Get(url)
		if err != nil {
			log.Printf("Error sending GET request: %v\n", err)
			continue
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("Error reading response body: %v\n", err)
			continue
		}

		var apiResponse model.GetFitMatesApiResponse
		if err := json.Unmarshal(body, &apiResponse); err != nil {
			log.Printf("Error unmarshalling API response: %v\n", err)
			continue
		}

		if err := fitMateService.HandleFitMateEvent(apiResponse); err != nil {
			log.Printf("Error handling fit mate event: %v\n", err)
		}
	}
}

func FitGroupHandler(c chan MessageEvent, fgService service.FitGroupServiceInterface) {
	for event := range c {
		msg := event.Message
		value, err := strconv.Atoi(string(msg.Value))
		if err != nil {
			log.Printf("Error converting Kafka message to int: %v\n", err)
			continue
		}
		url := fmt.Sprintf("http://fit-group:8080/fit-group-service/groups/%d", value)
		response, err := http.Get(url)
		if err != nil {
			log.Printf("Error sending GET request: %v\n", err)
			continue
		}
		defer response.Body.Close()

		var apiResponse model.GetFitGroupDetailApiResponse
		if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
			log.Printf("Error decoding API response: %v\n", err)
			continue
		}

		if err := fgService.HandleFitGroupEvent(apiResponse); err != nil {
			log.Printf("Error handling fit group event: %v\n", err)
		}
	}
}

func UserCreateEventHandler(c chan MessageEvent, userService service.UserService) {
	for event := range c {
		msg := event.Message
		value, err := strconv.Atoi(string(msg.Value))
		if err != nil {
			log.Printf("Error converting Kafka message to int: %v\n", err)
			continue
		}
		url := fmt.Sprintf("http://auth-service:8080/user/user-info/%d", value)
		response, err := http.Get(url)
		if err != nil {
			log.Printf("Error sending GET request: %v\n", err)
			continue
		}
		defer response.Body.Close()

		var apiResponse model.GetUserInfoApiResponse
		if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
			log.Printf("Error decoding API response: %v\n", err)
			continue
		}

		if err := userService.HandleUserCreateEvent(apiResponse); err != nil {
			log.Printf("Error handling user create event: %v\n", err)
		}
	}
}

func HandleMessage(msg kafka.Message, fitGroupService *service.FitGroupService, fitMateService service.FitMateService, userService service.UserService) {
	fitMateChannel := make(chan MessageEvent)
	fitGroupChannel := make(chan MessageEvent)
	userCreateEventChannel := make(chan MessageEvent)

	go FitMateHandler(fitMateChannel, fitMateService, fitGroupService)
	go FitGroupHandler(fitGroupChannel, fitGroupService)
	go UserCreateEventHandler(userCreateEventChannel, userService)

	switch msg.Topic {
	case "fit-mate":
		fitMateChannel <- MessageEvent{Message: msg, Service: fitMateService}
	case "fit-group":
		fitGroupChannel <- MessageEvent{Message: msg, Service: fitGroupService}
	case "user-create-event":
		userCreateEventChannel <- MessageEvent{Message: msg, Service: userService}
	default:
		fmt.Printf("No handler for topic %s\n", msg.Topic)
	}

	close(fitMateChannel)
	close(fitGroupChannel)
	close(userCreateEventChannel)
}
