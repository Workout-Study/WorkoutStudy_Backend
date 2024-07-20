package persistence

import (
	"database/sql"
	"log"
	"time"
	"workoutstudy_chatting/model"
)

type ChatRepository interface {
	RetrieveMessage(fitGroupID int) (int, error)
	RetrieveMessages(fitGroupID int, since time.Time) ([]model.ChatMessage, error)
	SaveMessage(msg model.ChatMessage) error
	RetrieveMessagesInRange(fitGroupID int, start, end time.Time) ([]model.ChatMessage, error)
}

type ChatRepositoryImpl struct {
	DB *sql.DB
}

// 훈기 tip : 인터페이스에 정의된 함수 중 구현안된거 체크
// var _ ChatRepository = &ChatRepositoryImpl{}
var _ ChatRepository = (*ChatRepositoryImpl)(nil)

func NewChatRepository(db *sql.DB) ChatRepository {
	return &ChatRepositoryImpl{DB: db}
}

// TODO : 서비스 레이어에서 이 함수를 사용해서 messageId 만 비교하도록 수정
func (repo *ChatRepositoryImpl) RetrieveMessage(fitGroupID int) (int, error) {
	query := `
	SELECT message_id FROM message WHERE fit_group_id = $1 ORDER BY message_time DESC LIMIT 1
	`
	log.Printf("Repository layer: Executing query for fitGroupID: %d", fitGroupID)
	var messageID int
	err := repo.DB.QueryRow(query, fitGroupID).Scan(&messageID)
	if err != nil {
		log.Printf("Repository layer: Error executing query: %v", err)
		return 0, err
	}
	return messageID, nil
}
func (repo *ChatRepositoryImpl) RetrieveMessages(fitGroupID int, since time.Time) ([]model.ChatMessage, error) {
	query := `
    SELECT message_id, user_id, fit_group_id, message, message_time, message_type
    FROM message
    WHERE fit_group_id = $1 AND message_time > $2
    ORDER BY message_time DESC
    `
	log.Printf("Repository layer: Executing query for fitGroupID: %d, since: %v", fitGroupID, since)
	rows, err := repo.DB.Query(query, fitGroupID, since)
	if err != nil {
		log.Printf("Repository layer: Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var messages []model.ChatMessage
	for rows.Next() {
		var msg model.ChatMessage
		if err := rows.Scan(&msg.ID, &msg.FitGroupID, &msg.Message, &msg.MessageTime, &msg.MessageType); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (repo *ChatRepositoryImpl) RetrieveMessagesInRange(fitGroupID int, start, end time.Time) ([]model.ChatMessage, error) {
	query := `
    SELECT message_id, user_id, fit_group_id, message, message_time, message_type
    FROM message
    WHERE fit_group_id = $1 AND message_time >= $2 AND message_time <= $3
    ORDER BY message_time ASC
    `
	log.Printf("Repository layer: Executing range query for fitGroupID: %d, start: %v, end: %v", fitGroupID, start, end)
	rows, err := repo.DB.Query(query, fitGroupID, start, end)
	if err != nil {
		log.Printf("Repository layer: Error executing range query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var messages []model.ChatMessage
	for rows.Next() {
		var msg model.ChatMessage
		if err := rows.Scan(&msg.ID, &msg.FitGroupID, &msg.Message, &msg.MessageTime, &msg.MessageType); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (repo *ChatRepositoryImpl) SaveMessage(msg model.ChatMessage) error {
	log.Printf("chat repository 에서 메시지 저장 시작: %v", msg)
	query := `
    INSERT INTO message (message_id, user_id, fit_group_id, message, message_time, message_type, created_at, created_by, updated_at, updated_by)
	VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7, NOW(), $7)
    `
	_, err := repo.DB.Exec(query, msg.ID, msg.FitGroupID, msg.UserID, msg.Message, msg.MessageTime, msg.MessageType, time.Now(), msg.UserID, time.Now(), msg.UserID)
	return err
}
