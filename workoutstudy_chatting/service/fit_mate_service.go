package service

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"workoutstudy_chatting/model"
	"workoutstudy_chatting/persistence"
)

type FitMateUseCase interface {
	GetFitGroupsByUserID(userID int) ([]model.FitGroup, error)
	GetFitMateByID(fitMateID string) (*model.FitMate, error)
	SaveFitMate(*model.FitMate) (*model.FitMate, error)
	DeleteFitMate(id int) ([]int, error)
	UpdateFitMate(*model.FitMate) (*model.FitMate, error)
	HandleFitMateEvent(apiResponse model.GetFitMatesApiResponse, fitGroupEvents chan int) error
}

// 인터페이스 구현 확인
var _ FitMateUseCase = (*FitMateService)(nil)

type FitMateService struct {
	repo           persistence.FitMateRepository
	fitGroupEvents chan int
}

func NewFitMateService(repo persistence.FitMateRepository, ch chan int) *FitMateService {
	return &FitMateService{
		repo:           repo,
		fitGroupEvents: ch,
	}
}

func (s *FitMateService) GetFitGroupsByUserID(fitMateID int) ([]model.FitGroup, error) {
	// 결과를 저장할 슬라이스 타입의 채널을 생성합니다.
	resultChan := make(chan []model.FitGroup)
	errorChan := make(chan error)

	go func() {
		// 저장소(repository)의 메서드를 호출합니다.
		fitGroups, err := s.repo.GetFitGroupsByUserID(fitMateID)
		if err != nil {
			// 에러가 발생하면 errorChan에 에러를 전송합니다.
			errorChan <- err
			return
		}
		// 에러가 없으면 resultChan에 결과를 전송합니다.
		resultChan <- fitGroups
	}()

	// select 문을 사용하여 결과 또는 에러를 기다립니다.
	select {
	case fitGroups := <-resultChan:
		// fitGroups를 성공적으로 받으면 반환합니다.
		return fitGroups, nil
	case err := <-errorChan:
		// 에러를 받으면 에러를 반환합니다.
		return nil, err
	}
}

func (s *FitMateService) GetFitMateByID(fitMateID string) (*model.FitMate, error) {
	// 결과를 저장할 포인터 타입의 채널을 생성합니다.
	resultChan := make(chan *model.FitMate)
	errorChan := make(chan error)

	go func() {
		// 저장소(repository)의 메서드를 호출합니다.
		fitMate, err := s.repo.GetFitMateByID(fitMateID)
		if err != nil {
			// 에러가 발생하면 errorChan에 에러를 전송합니다.
			errorChan <- err
			return
		}
		// 에러가 없으면 resultChan에 결과를 전송합니다.
		resultChan <- fitMate
	}()

	// select 문을 사용하여 결과 또는 에러를 기다립니다.
	select {
	case fitMate := <-resultChan:
		// FitMate 객체를 성공적으로 받으면 반환합니다.
		return fitMate, nil
	case err := <-errorChan:
		// 에러를 받으면 에러를 반환합니다.
		return nil, err
	}
}

func (s *FitMateService) SaveFitMate(fitMate *model.FitMate) (*model.FitMate, error) {
	return s.repo.SaveFitMate(fitMate)
}

func (s *FitMateService) DeleteFitMate(id int) ([]int, error) {
	return s.repo.DeleteFitMate(id)
}

func (s *FitMateService) UpdateFitMate(fitMate *model.FitMate) (*model.FitMate, error) {
	return s.repo.UpdateFitMate(fitMate)
}

/*
비교 및 조치 수행
1. Get Fit Mate list API 의 fitGroupId로 fit_group 테이블에서 fit_group 조회
1-a. fit_group 존재하지 않을 시 Wait -> fitGroup이 최초 생성되어 아직 Get Fit Group Detail API의 처리가 끝나지 않은 것.
1-b. fit_group 존재할 시 다음 단계 진행
2. fit_group 존재할 시 fitGroupId로 fit_mate 테이블 조회
2-a-1. 조회된 fit_mate 가 null -> 최초 생성된 fit_group 임을 의미
2-a-2. 조회된 fit_mate 가 null 이 아닐 시 다음 단계 진행
3. Response 의 Mate 정보와 fit_mate 조회 결과 비교
4. Reponse 의 Mate 대로 fit_mate UPDATE & DELETE
4-a. UPDATE : Response 애는 존재하고 fit_mate에는 없는 경우
4-a-1. Response의 Mate 대로 fit_mate 생성(INSERT
4-b. DELETE : Response에는 없고 fit_mate에는 존재하는 경우
4-b-1. fit_mate 삭제(DELETE), Hard Delete 로 진행
*/
// fit_mate_service.go
func (s *FitMateService) HandleFitMateEvent(apiResponse model.GetFitMatesApiResponse, fitGroupEvents chan int) error {
	// fitGroupId로 fit_group 테이블 조회
	fitGroupExists, err := s.repo.CheckFitGroupExists(apiResponse.FitGroupId)
	if err != nil {
		log.Printf("Error checking fit group existence: %v", err)
		return err
	}

	// fit_group 존재하지 않을 시 10초 대기
	if !fitGroupExists {
		log.Printf("FitGroup ID %d does not exist. Waiting...", apiResponse.FitGroupId)
		timeout := time.After(10 * time.Second)
		select {
		case fitGroupID := <-fitGroupEvents:
			if fitGroupID != apiResponse.FitGroupId {
				log.Printf("Waiting for FitGroup ID %d, but received %d", apiResponse.FitGroupId, fitGroupID)
				return nil // 혹은 적절한 에러 처리
			}
			log.Printf("Proceeding with FitGroup ID %d", fitGroupID)
		case <-timeout:
			log.Printf("Timeout waiting for FitGroup ID %d", apiResponse.FitGroupId)
			// 10초 후 DB 재조회
			fitGroupExists, err = s.repo.CheckFitGroupExists(apiResponse.FitGroupId)
			if err != nil {
				log.Printf("Error checking fit group existence after timeout: %v", err)
				return err
			}
			if !fitGroupExists {
				log.Printf("FitGroup ID %d not found after 10 seconds", apiResponse.FitGroupId)
				return fmt.Errorf("fitGroup ID %d not found after 10 seconds", apiResponse.FitGroupId)
			}
		}
	}

	fitMateIds, err := s.repo.GetFitMatesIdsByFitGroupId(apiResponse.FitGroupId)
	if err != nil {
		log.Printf("Error fetching fit mate IDs from DB: %v", err)
		return err
	}

	return s.compareAndUpdateFitMates(apiResponse, fitMateIds)
}

func (s *FitMateService) compareAndUpdateFitMates(apiResponse model.GetFitMatesApiResponse, dbFitMateIds []int) error {
	apiFitMateIdsMap := make(map[int]bool)
	dbFitMateMap := make(map[int]*model.FitMate)

	// FitMateDetails 처리 로직
	for _, detail := range apiResponse.FitMateDetails {
		apiFitMateIdsMap[detail.FitMateId] = true
		if dbFitMate, err := s.repo.GetFitMateByID(strconv.Itoa(detail.FitMateId)); err == nil {
			dbFitMateMap[detail.FitMateId] = dbFitMate
		}
	}

	// DB에 존재하는 FitMate들 삭제
	for _, dbId := range dbFitMateIds {
		if !apiFitMateIdsMap[dbId] {
			_, err := s.repo.DeleteFitMate(dbId)
			if err != nil {
				log.Printf("Error deleting fit mate ID %d: %v", dbId, err)
				return err
			}
		}
	}

	// 새로운 FitMate들 추가
	for _, apiDetail := range apiResponse.FitMateDetails {
		if _, exists := dbFitMateMap[apiDetail.FitMateId]; !exists {
			newFitMate := &model.FitMate{
				ID:         apiDetail.FitMateId,
				UserID:     apiDetail.FitMateUserId, // Change the field name from "UserID" to "UserID"
				FitGroupID: apiResponse.FitGroupId,
				State:      false,
				CreatedBy:  "system",
			}
			_, err := s.repo.SaveFitMate(newFitMate)
			if err != nil {
				log.Printf("Error adding new fit mate ID %d: %v", apiDetail.FitMateId, err)
				return err
			}
		}
	}

	return nil
}

// TODO : fit mate UPDATE 할 게 현재 딱히 없음. state 는 조회 결과에 없고, nickname 은 user-create-event로 처리함
// for _, apiDetail := range apiResponse.FitMateDetails {
// 	if dbFitMate, exists := dbFitMateMap[apiDetail.FitMateId]; exists {
// 		if dbFitMate.State != apiDetail.State {
// 			_, err := s.UpdateFitMate(dbFitMate)
// 			if err != nil {
// 				log.Printf("Error updating fit mate ID %d: %v", apiDetail.FitMateId, err)
// 				return err
// 			}
// 		}
// 	}
// }
