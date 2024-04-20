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

// fit-mate 이벤트 컨슘 후 호출하는 API 응답 타입 정의
type GetFitMatesApiResponse struct {
	FitMateDetails []struct {
		FitMateId int `json:"fitMateid"`
	} `json:"fitMateDetails"`
}

/*
	TODO :
	1. DB 싱크를 정일님과 맞춰야함
	2. API Body 결과의 fitMateDetails[].fitMateId 와
	DB 에서 fit_group_mate 테이블에서 fit-group-id 로 조회해서 나온
	리스트와 비교해서
	1. 추가 : fitMateDetails[].fitMateId 에는 존재하는데 DB 리스트에는 없으면 DB 에 없는 사용자 추가
	2. 삭제 : fitMateDetails[].fitMateId 에는 존재하지 않는데 DB 리스트에는 있으면 DB 에서 삭제
	3. 보존 : fitMateDetails[].fitMateId 와 DB 차이가 없으면 그대로 보존
*/

// FitMateHandler는 'fit-mate' 토픽의 메시지를 컨슘
func FitMateHandler(msg kafka.Message, fitGroupService *service.FitGroupService, fitMateService service.FitMateService) {
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

	var apiResponse GetFitMatesApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("Error unmarshalling API response: %v\n", err)
		return
	}

	dbFitMateIds, err := fitGroupService.GetFitMatesByFitGroupId(value)
	if err != nil {
		log.Printf("Error fetching fit mate IDs from DB: %v\n", err)
		return
	}

	// API 결과와 DB 결과 비교 및 조치 수행
	compareFitMateDetails(apiResponse, dbFitMateIds, fitMateService)

	fmt.Printf("Response from fit-group-service for groups: %s\n", string(body))
}

// compareFitMateDetails는 API 결과와 DB 결과를 비교하여 필요한 조치를 수행합니다.
func compareFitMateDetails(apiResponse GetFitMatesApiResponse, dbFitMateIds []int, fitMateService service.FitMateService) {
	apiFitMateIdsMap := make(map[int]bool)
	for _, detail := range apiResponse.FitMateDetails {
		apiFitMateIdsMap[detail.FitMateId] = true
	}

	// DB에 있지만 API에는 없는 경우 - 삭제
	for _, dbId := range dbFitMateIds {
		if !apiFitMateIdsMap[dbId] {
			log.Printf("Fit mate ID %d is in DB but not in API response, should be deleted\n", dbId)
		}
	}

	// API에는 있지만 DB에는 없는 경우 - 추가
	for _, apiDetail := range apiResponse.FitMateDetails {
		if _, exists := apiFitMateIdsMap[apiDetail.FitMateId]; !exists {
			newFitMate := &model.FitMate{ // 이미 포인터로 선언되어 있으므로 이 부분은 변경 없음
				ID:        apiDetail.FitMateId,
				Username:  "새 유저",
				Nickname:  "닉네임",
				State:     true,
				CreatedBy: "system",
			}
			_, err := fitMateService.SaveFitMate(newFitMate) // 이제 인터페이스에 맞춰 올바르게 호출
			if err != nil {
				log.Printf("Error adding new fit mate ID %d: %v", apiDetail.FitMateId, err)
			} else {
				log.Printf("Added new fit mate ID %d to the database", apiDetail.FitMateId)
			}
		}
	}
}

// FitGroupHandler는 'fit-group' 토픽의 메시지를 컨슘
func FitGroupHandler(msg kafka.Message) {
	value, err := strconv.Atoi(string(msg.Value))
	log.Printf("kafka fit-group event value: %d\n", value)
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
	/*
		TODO
		Register, Update, Delete 시에 이벤트 들어온
		조회 응답 바디가 DB에서 조회한 내용과 차이가 나면 update
		최초 Register 시에는 내 쪽 DB에서 조회하면 null 나올 테니
		그러면 생성
	*/
	fmt.Printf("Response from fit-group-service for mates: %s\n", string(body))
}

// HandleMessage는 토픽에 따라 적절한 핸들러를 호출합니다.
func HandleMessage(msg kafka.Message, fitGroupService *service.FitGroupService, fitMateService service.FitMateService) {
	switch msg.Topic {
	case "fit-mate":
		FitMateHandler(msg, fitGroupService, fitMateService)
	case "fit-group":
		FitGroupHandler(msg)
	default:
		fmt.Printf("No handler for topic %s\n", msg.Topic)
	}
}
