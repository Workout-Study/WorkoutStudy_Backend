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
	// 로그로 since 파라미터 출력
	log.Printf("Repository layer: 최신 채팅 조회 시작 fitGroupID: %d, since: %v", fitGroupID, since)

	query := `
    SELECT message_id, user_id, fit_group_id, message, message_time, message_type
    FROM message
    WHERE fit_group_id = $1 AND message_time >= $2 -- 최신 메시지와 요청 메시지 시간 포함
    ORDER BY message_time DESC
    LIMIT 1 -- 최신 메시지 하나만 가져옴
    `
	rows, err := repo.DB.Query(query, fitGroupID, since)
	if err != nil {
		log.Printf("Repository layer: Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var messages []model.ChatMessage
	for rows.Next() {
		var msg model.ChatMessage
		if err := rows.Scan(&msg.ID, &msg.UserID, &msg.FitGroupID, &msg.Message, &msg.MessageTime, &msg.MessageType); err != nil {
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
	log.Printf("Repository layer: 채팅 범위 조회 시작 fitGroupID: %d, start: %v, end: %v", fitGroupID, start, end)
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
		if err := rows.Scan(&msg.ID, &msg.UserID, &msg.FitGroupID, &msg.Message, &msg.MessageTime, &msg.MessageType); err != nil {
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
	log.Printf("chat repository 에서 저장하는 메세지의 파싱 전 messageTime: %v", msg.MessageTime)
	formattedTime := msg.MessageTime.Format("2006-01-02 15:04:05.999999-07:00")
	log.Printf("chat repository 에서 저장하는 메세지의 파싱 후 messageTime: %v", formattedTime)
	query := `
    INSERT INTO message (message_id, user_id, fit_group_id, message, message_time, message_type, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
    `
	_, err := repo.DB.Exec(query, msg.ID, msg.UserID, msg.FitGroupID, msg.Message, formattedTime, msg.MessageType)
	return err
}
