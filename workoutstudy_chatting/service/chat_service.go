package service

import (
	"log"
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

func (service *ChatService) RetrieveMessages(fitGroupID int, since time.Time, messageID string) ([]model.ChatMessage, string, error) {
	log.Printf("Service layer: Retrieving messages for fitGroupID: %d, since: %v", fitGroupID, since)
	messages, err := service.repo.RetrieveMessages(fitGroupID, since)
	if err != nil {
		return nil, "", err
	}
	log.Printf("Service layer: Retrieved messages count: %d", len(messages))

	var filteredMessages []model.ChatMessage
	var latestMessageId string
	if len(messages) > 0 {
		latestMessageId = messages[0].ID // 최신 메시지 ID 저장
		if latestMessageId == messageID {
			// 두 messageId 가 일치할 시 최신 message 객체만 Return
			filteredMessages = []model.ChatMessage{messages[0]}
		} else {
			// 두 messageId가 불일치할 시 모든 message 객체를 Return
			filteredMessages = messages
		}
	}

	// 필터링된 메시지 배열과 최신 메시지의 ID 반환
	return filteredMessages, latestMessageId, nil
}

func (service *ChatService) SaveChatMessage(msg model.ChatMessage) error {
	return service.repo.SaveMessage(msg)
}
