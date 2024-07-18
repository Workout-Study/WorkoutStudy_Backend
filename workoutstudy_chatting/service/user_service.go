package service

import (
	"fmt"
	"log"
	"workoutstudy_chatting/model"
	"workoutstudy_chatting/persistence"
)

type UserUseCase interface {
	SaveUser(user *model.User) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	DeleteUser(userID int) error
	GetUserByID(userID int) (*model.User, error)
	HandleUserCreateEvent(user *model.UserCreateEvent) error
	HandleUserInfoEvent(apiResponse model.GetUserInfoApiResponse) error
}

// 컴파일 타임에 인터페이스 구현 확인
var _ UserUseCase = (*UserService)(nil)

/*
var _ UserUseCase = &UserService{} 이 코드는 UserService 인스턴스를 실제로 생성하여
UserUseCase 인터페이스에 할당하는 코드. 생성자 코드가 있다면 중복되는 것.
*/

type UserService struct {
	repo persistence.UserRepository
}

/*
* : 포인터는 다른 변수의 메모리 주소를 저장하는 변수. '*'는 포인터 타입을 정의할 때 사용
& : 주소 연산자. 변수의 메모리 주소를 반환
*/
func NewUserService(repo persistence.UserRepository) *UserService {
	return &UserService{repo: repo}
	// return type이 UserService 타입의 포인터로 라는 뜻. *UserService는 UserService 타입의 포인터를 의미
	// 실제 return 문에서 UserService 의 메모리 주소를 반환
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
func (s *UserService) HandleUserCreateEvent(userCreateEvent *model.UserCreateEvent) error {
	user := &model.User{
		ID:        userCreateEvent.UserID,
		Nickname:  userCreateEvent.Nickname,
		State:     userCreateEvent.State,
		CreatedBy: userCreateEvent.Nickname,
		UpdatedBy: userCreateEvent.Nickname,
	}

	_, err := s.repo.SaveUser(user)
	if err != nil {
		log.Printf("Error handling user create event: %v", err)
		return fmt.Errorf("error handling user create event: %w", err)
	}

	return nil
}

// Update
func (s *UserService) HandleUserInfoEvent(apiResponse model.GetUserInfoApiResponse) error {
	// GetUserInfoApiResponse를 User 모델로 변환
	user := &model.User{
		ID:       apiResponse.UserID,
		Nickname: apiResponse.Nickname,
		// State:     false,                // 상태를 true로 설정하거나 필요한 로직에 따라 변경하세요
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
