package service

import (
	"workoutstudy_chatting/model"
	"workoutstudy_chatting/persistence"
)

type FitMateService interface {
	GetFitGroupByMateID(fitMateID string) ([]model.FitGroup, error)
	GetFitMateByID(fitMateID string) (*model.FitMate, error)
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
