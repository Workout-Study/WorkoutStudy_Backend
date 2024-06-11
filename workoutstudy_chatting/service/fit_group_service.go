package service

import (
	"log"
	"workoutstudy_chatting/model"
	"workoutstudy_chatting/persistence"
)

type FitGroupUseCase interface {
	GetFitGroupByID(fitGroupID int) (*model.FitGroup, error)
	GetFitMatesByFitGroupId(fitGroupID int) ([]int, error)
	SaveFitGroup(fitGroup *model.FitGroup) (*model.FitGroup, error)
	HandleFitGroupEvent(apiResponse model.GetFitGroupDetailApiResponse, fitGroupEvents chan int) error
}

// 인터페이스 구현 확인
var _ FitGroupUseCase = (*FitGroupService)(nil)

type FitGroupService struct {
	repo            persistence.FitGroupRepository
	fitGroupCreated chan int
}

func NewFitGroupService(repo persistence.FitGroupRepository, ch chan int) *FitGroupService {
	return &FitGroupService{repo: repo, fitGroupCreated: ch}
}

func (s *FitGroupService) GetFitGroupByID(fitGroupID int) (*model.FitGroup, error) {
	fitGroup, err := s.repo.GetFitGroupByID(fitGroupID)
	if err != nil {
		return nil, err
	}
	return fitGroup, nil
}

func (s *FitGroupService) GetFitMatesByFitGroupId(fitGroupID int) ([]int, error) {
	fitMateIds, err := s.repo.GetFitMatesByFitGroupId(fitGroupID)
	if err != nil {
		return nil, err
	}
	return fitMateIds, nil
}

func (s *FitGroupService) SaveFitGroup(fitGroup *model.FitGroup) (*model.FitGroup, error) {
	return s.repo.SaveFitGroup(fitGroup)
}

func (s *FitGroupService) HandleFitGroupEvent(apiResponse model.GetFitGroupDetailApiResponse, fitGroupEvents chan int) error {
	// 1. Get Fit group detail API 의 fitGroupId로 fit_group 테이블 조회
	fitGroup, err := s.repo.GetFitGroupByID(apiResponse.FitGroupId)
	if err != nil {
		return err
	}

	// 1-a. 존재할 시 Create skip -> Delete 로 이동
	if fitGroup != nil {
		// 2. API Response 의 state 확인
		if apiResponse.State {
			// 2-a. state 가 true 일 시 DB에서 해당 fit_group의 state 를 true 로 변경 -> 삭제
			return s.repo.DeleteFitGroup(apiResponse.FitGroupId)
		} else {
			// 2-b. state 가 false 일 시 진행 skip, proceed to Update
			// Update
			// 2. API Response 와 DB 의 fit_group 정보 비교
			if shouldUpdate(fitGroup, apiResponse) {
				// 2-a. 다를 시 DB 정보를 API Response 로 업데이트
				return s.repo.UpdateFitGroup(convertApiToModel(apiResponse))
			}
			// 2-b. 같을 시 진행 skip
			return nil // No changes needed, nothing to update
		}
	} else {
		// 1-b-1. DB에 존재하지 않을 시 fit_group 테이블에 API Response 로 row 생성
		newFitGroup, err := s.repo.SaveFitGroup(convertApiToModel(apiResponse))
		if err != nil {
			return err
		}
		// 1-b-2. row 생성 이후, fit_group row가 입력됐다는 것을 Get Fit Mate list API Handler 에게 알림
		select {
		case fitGroupEvents <- newFitGroup.ID:
			log.Printf("Notified fit mate handler of new fit group ID %d", newFitGroup.ID)
		default:
			log.Printf("Failed to notify fit mate handler of new fit group ID %d", newFitGroup.ID)
		}
		return nil
	}
}

func shouldUpdate(existing *model.FitGroup, response model.GetFitGroupDetailApiResponse) bool {
	// Example of comparison logic; extend this based on actual fields that matter
	return existing.FitLeaderUserID != response.FitLeaderUserId ||
		existing.FitGroupName != response.FitGroupName ||
		existing.State != response.State
}

func convertApiToModel(apiResp model.GetFitGroupDetailApiResponse) *model.FitGroup {
	return &model.FitGroup{
		ID:                  apiResp.FitGroupId,
		FitLeaderUserID:     apiResp.FitLeaderUserId,
		FitGroupName:        apiResp.FitGroupName,
		Category:            apiResp.Category,
		Cycle:               apiResp.Cycle,
		Frequency:           apiResp.Frequency,
		PresentFitMateCount: apiResp.PresentFitMateCount,
		MaxFitMate:          apiResp.MaxFitMate,
		State:               apiResp.State,
	}
}
