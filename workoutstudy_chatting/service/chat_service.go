package service

import (
	"time"
	"workoutstudy_chatting/model"
	"workoutstudy_chatting/persistence"
)

type ChatService struct {
	repo persistence.ChatRepository
}

func NewChatService(repo persistence.ChatRepository) *ChatService {
	return &ChatService{repo: repo}
}

func (service *ChatService) RetrieveMessages(fitGroupID int, since time.Time) ([]model.ChatMessage, error) {
	// 비동기 처리를 위한 채널 선언
	resultChan := make(chan []model.ChatMessage)
	errorChan := make(chan error)

	go func() {
		messages, err := service.repo.RetrieveMessages(fitGroupID, since)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- messages
	}()

	select {
	case messages := <-resultChan:
		return messages, nil
	case err := <-errorChan:
		return nil, err
	}
}

func (service *ChatService) SaveChatMessage(msg model.ChatMessage) error {
	return service.repo.SaveMessage(msg)
}
