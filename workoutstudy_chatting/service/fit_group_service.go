package service

import (
	"workoutstudy_chatting/model"
	"workoutstudy_chatting/persistence"
)

type FitGroupServiceInterface interface {
	GetFitGroupByID(fitGroupID int) (*model.FitGroup, error)
}

type FitGroupService struct {
	repo persistence.FitGroupRepository
}

func NewFitGroupService(repo persistence.FitGroupRepository) *FitGroupService {
	return &FitGroupService{repo: repo}
}

func (s *FitGroupService) GetFitGroupByID(fitGroupID int) (bool, error) {
	fitGroup, err := s.repo.GetFitGroupByID(fitGroupID)
	if err != nil {
		return false, err
	}
	return fitGroup != nil, nil
}
