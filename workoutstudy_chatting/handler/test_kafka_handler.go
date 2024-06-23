package handler

import (
	"encoding/json"
	"log"
	"strconv"

	"workoutstudy_chatting/model"
	"workoutstudy_chatting/service"

	"github.com/segmentio/kafka-go"
)

func HandleFitMateEvent(msg kafka.Message, fitMateService service.FitMateUseCase) {
	value, err := strconv.Atoi(string(msg.Value))
	if err != nil {
		log.Printf("Error converting Kafka message to int: %v\n", err)
		return
	}

	// 여기서는 간단히 로그를 출력합니다. 실제 로직을 호출할 수 있습니다.
	log.Printf("Handling fit-mate event with value: %d\n", value)
	// 실제 서비스 로직 호출
	// fitMateService.SaveFitMate(value)
}

func HandleFitGroupEvent(msg kafka.Message, fitGroupService service.FitGroupUseCase) {
	value, err := strconv.Atoi(string(msg.Value))
	if err != nil {
		log.Printf("Error converting Kafka message to int: %v\n", err)
		return
	}

	// 여기서는 간단히 로그를 출력합니다. 실제 로직을 호출할 수 있습니다.
	log.Printf("Handling fit-group event with value: %d\n", value)
	// 실제 서비스 로직 호출
	// fitGroupService.HandleFitGroup(value)
}

func HandleUserCreateEvent(msg kafka.Message, userService service.UserUseCase) {
	var userCreateEvent model.UserCreateEvent

	if err := json.Unmarshal(msg.Value, &userCreateEvent); err != nil {
		log.Printf("Error unmarshalling message: %v\n", err)
		return
	}

	// 여기서는 간단히 로그를 출력합니다. 실제 로직을 호출할 수 있습니다.
	log.Printf("Handling user-create event with user: %v\n", userCreateEvent)
	// 실제 서비스 로직 호출
	// userService.HandleUserCreate(&userCreateEvent)
}

func HandleUserInfoEvent(msg kafka.Message, userService service.UserUseCase) {
	value, err := strconv.Atoi(string(msg.Value))
	if err != nil {
		log.Printf("Error converting Kafka message to int: %v\n", err)
		return
	}

	// 여기서는 간단히 로그를 출력합니다. 실제 로직을 호출할 수 있습니다.
	log.Printf("Handling user-info event with value: %d\n", value)
	// 실제 서비스 로직 호출
	// userService.HandleUserInfo(value)
}
