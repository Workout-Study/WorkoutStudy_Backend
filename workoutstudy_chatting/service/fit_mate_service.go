package service

import (
	"log"
	"workoutstudy_chatting/model"
	"workoutstudy_chatting/persistence"
)

type FitMateService interface {
	GetFitGroupByMateID(fitMateID string) ([]model.FitGroup, error)
	GetFitMateByID(fitMateID string) (*model.FitMate, error)
	SaveFitMate(*model.FitMate) (*model.FitMate, error)
	DeleteFitMate(id int) ([]int, error)
	HandleFitMateEvent(apiResponse model.GetFitMatesApiResponse) error
}

type FitMateServiceImpl struct {
	repo persistence.FitMateRepository
}

var _ FitMateService = &FitMateServiceImpl{}

func NewFitMateService(repo persistence.FitMateRepository) FitMateService {
	return &FitMateServiceImpl{repo: repo}
}

func (s *FitMateServiceImpl) GetFitGroupByMateID(fitMateID string) ([]model.FitGroup, error) {
	// 결과를 저장할 슬라이스 타입의 채널을 생성합니다.
	resultChan := make(chan []model.FitGroup)
	errorChan := make(chan error)

	go func() {
		// 저장소(repository)의 메서드를 호출합니다.
		fitGroups, err := s.repo.GetFitGroupByMateID(fitMateID)
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

func (s *FitMateServiceImpl) GetFitMateByID(fitMateID string) (*model.FitMate, error) {
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

func (s *FitMateServiceImpl) SaveFitMate(fitMate *model.FitMate) (*model.FitMate, error) {
	return s.repo.SaveFitMate(fitMate)
}

func (s *FitMateServiceImpl) DeleteFitMate(id int) ([]int, error) {
	return s.repo.DeleteFitMate(id)
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
func (s *FitMateServiceImpl) HandleFitMateEvent(apiResponse model.GetFitMatesApiResponse) error {
	fitMateIds, err := s.repo.GetFitMatesIdsByFitGroupId(apiResponse.FitGroupId)
	if err != nil {
		log.Printf("Error fetching fit mate IDs from DB: %v\n", err)
		return err
	}

	// API 결과와 DB 결과 비교 및 조치 수행
	return s.compareAndUpdateFitMates(apiResponse, fitMateIds)
}

func (s *FitMateServiceImpl) compareAndUpdateFitMates(apiResponse model.GetFitMatesApiResponse, dbFitMateIds []int) error {
	apiFitMateIdsMap := make(map[int]bool)
	for _, detail := range apiResponse.FitMateDetails {
		apiFitMateIdsMap[detail.FitMateId] = true
	}

	// DB에 있지만 API에는 없는 경우 - 삭제
	for _, dbId := range dbFitMateIds {
		if _, found := apiFitMateIdsMap[dbId]; !found {
			if _, err := s.DeleteFitMate(dbId); err != nil {
				log.Printf("Error deleting fit mate ID %d: %v", dbId, err)
				return err
			}
		}
	}

	// API에는 있지만 DB에는 없는 경우 - 추가
	for _, apiDetail := range apiResponse.FitMateDetails {
		if _, exists := apiFitMateIdsMap[apiDetail.FitMateId]; !exists {
			newFitMate := &model.FitMate{
				ID:        apiDetail.FitMateId,
				Username:  "새 유저",
				Nickname:  "닉네임",
				State:     true,
				CreatedBy: "system",
			}
			if _, err := s.SaveFitMate(newFitMate); err != nil {
				log.Printf("Error adding new fit mate ID %d: %v", apiDetail.FitMateId, err)
				return err
			}
		}
	}

	return nil
}
