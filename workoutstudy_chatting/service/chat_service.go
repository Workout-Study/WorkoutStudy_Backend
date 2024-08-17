package service

import (
	"log"
	"time"
	"workoutstudy_chatting/model"
	"workoutstudy_chatting/persistence"
)

type ChatUseCase interface {
	RetrieveMessages(fitGroupID int, messageTime time.Time, messageID string) ([]model.ChatMessage, string, error)
	SaveChatMessage(msg model.ChatMessage) error
}

var _ ChatUseCase = (*ChatService)(nil)

type ChatService struct {
	repo persistence.ChatRepository
}

func NewChatService(repo persistence.ChatRepository) *ChatService {
	return &ChatService{repo: repo}
}

/*
RetrieveMesaages
1. message 테이블에서 최신 채팅 조회
2. 조회된 message 의 messageId와 요청의 messageId 비교
3-a. 두 messageId 가 일치할 시 최신 채팅만 반환
3-b. 두 messageId 가 불일치할 시 요청 messageId 와 조회된 최신 채팅 사이의 채팅을 반환
*/
func (s *ChatService) RetrieveMessages(fitGroupID int, messageTime time.Time, messageID string) ([]model.ChatMessage, string, error) {
	log.Printf("Service layer: Retrieving messages for fit_group_id: %d, messageTime: %v", fitGroupID, messageTime)
	log.Printf("최신 메시지 조회 시작")

	// 최신 메시지들을 조회
	messages, err := s.repo.RetrieveMessages(fitGroupID, messageTime)
	if err != nil {
		return nil, "", err
	}
	log.Printf("Service layer: Retrieved messages count: %d", len(messages))

	var filteredMessages []model.ChatMessage
	var latestMessageId string

	// 최신 메시지가 존재하는 경우
	if len(messages) > 0 {
		latestMessageId = messages[0].ID // 최신 메시지 ID 저장

		// 최신 메시지와 요청된 messageID 비교
		if latestMessageId == messageID {
			// 두 messageID가 일치할 시 최신 메시지만 반환
			filteredMessages = []model.ChatMessage{messages[0]}
		} else {
			// 메시지 ID가 일치하지 않으면 요청된 시간 이후 메시지들 반환
			filteredMessages, err = s.repo.RetrieveMessagesInRange(fitGroupID, messageTime, messages[0].MessageTime)
			if err != nil {
				return nil, "", err
			}
		}
	} else {
		// 메시지가 없는 경우 빈 리스트 반환
		return []model.ChatMessage{}, "", nil
	}

	return filteredMessages, latestMessageId, nil
}

func (s *ChatService) SaveChatMessage(msg model.ChatMessage) error {
	return s.repo.SaveMessage(msg)
}
