package service

import (
	"fmt"
	"log"
	"workoutstudy_chatting/model"
	"workoutstudy_chatting/persistence"
)

type UserServiceInterface interface {
	SaveUser(user *model.User) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	DeleteUser(userID int) error
	GetUserByID(userID int) (*model.User, error)
	HandleUserCreateEvent(apiResponse model.GetUserInfoApiResponse) error
}

type UserService struct {
	repo persistence.UserRepository
}

var _ UserServiceInterface = &UserService{}

func NewUserService(repo persistence.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) SaveUser(user *model.User) (*model.User, error) {
	return s.repo.SaveUser(user)
}

func (s *UserService) UpdateUser(user *model.User) (*model.User, error) {
	return s.repo.UpdateUser(user)
}

func (s *UserService) DeleteUser(userID int) error {
	return s.repo.DeleteUser(userID)
}

func (s *UserService) GetUserByID(userID int) (*model.User, error) {
	return s.repo.GetUserByID(userID)
}

// Create
func (s *UserService) HandleUserCreateEvent(apiResponse model.GetUserInfoApiResponse) error {
	// GetUserInfoApiResponse를 User 모델로 변환
	user := &model.User{
		ID:        apiResponse.UserID,
		Nickname:  apiResponse.Nickname,
		State:     false,                // 상태를 true로 설정하거나 필요한 로직에 따라 변경하세요
		CreatedBy: apiResponse.Nickname, // 생성자를 Nickname으로 설정하거나 필요한 로직에 따라 변경하세요
		UpdatedBy: apiResponse.Nickname, // 업데이트한 사람을 Nickname으로 설정하거나 필요한 로직에 따라 변경하세요
	}

	// SaveUser 함수를 호출하여 사용자를 저장
	_, err := s.repo.SaveUser(user)
	if err != nil {
		log.Printf("Error handling user create event: %v", err)
		return fmt.Errorf("error handling user create event: %w", err)
	}

	return nil
}
