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

func HandleMessage(msgChan chan MessageEvent, fitMateService service.FitMateUseCase, fitGroupService service.FitGroupUseCase, userService service.UserUseCase) {
	fitMateChannel := make(chan MessageEvent)
	fitGroupChannel := make(chan MessageEvent)
	userCreateEventChannel := make(chan MessageEvent)
	userInfoEventChannel := make(chan MessageEvent)

	// fitGroupEvents 채널 생성
	fitGroupEvents := make(chan int, 1) // 비동기 이벤트 알림을 위해 채널 사용

	go FitMateHandler(fitMateChannel, fitMateService, fitGroupEvents)
	go FitGroupHandler(fitGroupChannel, fitGroupService, fitGroupEvents)
	go UserCreateEventHandler(userCreateEventChannel, userService)
	go UserInfoHandler(userInfoEventChannel, userService)

	for msgEvent := range msgChan {
		msg := msgEvent.Message

		switch msg.Topic {
		case "fit-mate":
			log.Printf("fit-mate 이벤트 컨슘: %s", string(msg.Value))
			fitMateChannel <- msgEvent
		case "fit-group":
			log.Printf("fit-group 이벤트 컨슘: %s", string(msg.Value))
			fitGroupChannel <- msgEvent
		case "user-create-event":
			log.Printf("user-create-event 이벤트 컨슘: %s", string(msg.Value))
			userCreateEventChannel <- msgEvent
		case "user-info":
			log.Printf("user-info 이벤트 컨슘: %s", string(msg.Value))
			userInfoEventChannel <- msgEvent
		default:
			fmt.Printf("No handler for topic %s\n", msg.Topic)
		}
	}
}

func FitMateHandler(c chan MessageEvent, fitMateService service.FitMateUseCase, fitGroupEvents chan int) {
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

		func() {
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

			if err := fitMateService.HandleFitMateEvent(apiResponse, fitGroupEvents); err != nil {
				log.Printf("Error handling fit mate event: %v\n", err)
			}
		}()
	}
}

func FitGroupHandler(c chan MessageEvent, fgService service.FitGroupUseCase, fitGroupEvents chan int) {
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

		// 익명함수를 사용해서 response.Body.Close()를 defer로 사용하지 않고 명시적으로 호출
		// 익명함수를 선언하고 그 안에서 defer reseponse.Body.Close()를 호출
		// 이렇게 하면 response.Body.Close()가 호출되는 시점을 명확히 알 수 있음
		// API를 호출해서 데이터를 받아온 후, 자원 해제 -> 그 후에 데이터에 대한 처리 수행
		func() {
			defer response.Body.Close()

			var apiResponse model.GetFitGroupDetailApiResponse
			if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
				log.Printf("Error decoding API response: %v\n", err)
				return
			}

			if err := fgService.HandleFitGroupEvent(apiResponse, fitGroupEvents); err != nil {
				log.Printf("Error handling fit group event: %v\n", err)
			}
		}()
	}
}

func UserCreateEventHandler(c chan MessageEvent, userService service.UserUseCase) {
	for event := range c {
		msg := event.Message

		var userCreateEvent model.UserCreateEvent

		if err := json.Unmarshal(msg.Value, &userCreateEvent); err != nil {
			log.Printf("Error unmarshalling message: %v\n", err)
			continue
		}

		if err := userService.HandleUserCreateEvent(&userCreateEvent); err != nil {
			log.Printf("Error handling user creation process: %v\n", err)
		}
	}
}

func UserInfoHandler(c chan MessageEvent, userService service.UserUseCase) {
	for event := range c {
		msg := event.Message
		value, err := strconv.Atoi(string(msg.Value))
		if err != nil {
			log.Printf("Error converting Kafka message to int: %v\n", err)
			continue
		}
		url := fmt.Sprintf("http://auth-service:8080/user/user-info?userId=%d", value)
		response, err := http.Get(url)
		if err != nil {
			log.Printf("Error sending GET request: %v\n", err)
			continue
		}

		func() {
			defer response.Body.Close()

			var apiResponse model.GetUserInfoApiResponse
			if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
				log.Printf("Error decoding API response: %v\n", err)
				return
			}

			if err := userService.HandleUserInfoEvent(apiResponse); err != nil {
				log.Printf("Error handling user info event: %v\n", err)
			}
		}()
	}
}
