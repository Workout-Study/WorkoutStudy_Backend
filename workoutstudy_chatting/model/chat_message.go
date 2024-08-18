package model

import (
	"encoding/json"
	"log"
	"time"
)

// MessageType은 메시지 타입을 나타내는 사용자 정의 타입입니다.
type MessageType string

// 가능한 MessageType 값을 상수로 정의합니다.
const (
	Chatting MessageType = "CHATTING"
	Ticket   MessageType = "TICKET"
)

// ChatMessage는 채팅 메시지를 나타내는 구조체입니다.
type ChatMessage struct {
	ID          string      `json:"messageId"`
	UserID      int         `json:"userId"`
	FitGroupID  int         `json:"fitGroupId"`
	Message     string      `json:"message"`
	MessageTime time.Time   `json:"messageTime"`
	MessageType MessageType `json:"messageType"`
}

// UnmarshalJSON은 시간 문자열을 time.Time으로 변환할 때 타임존 정보가 포함되도록 파싱합니다.
func (cm *ChatMessage) UnmarshalJSON(data []byte) error {
	type Alias ChatMessage
	tmp := struct {
		MessageTime string `json:"messageTime"`
		*Alias
	}{
		Alias: (*Alias)(cm),
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	// 타임존 정보를 포함한 시간 문자열을 파싱
	t, err := time.Parse(time.RFC3339Nano, tmp.MessageTime)
	if err != nil {
		log.Printf("Error parsing MessageTime: %v", err)
		return err
	}

	cm.MessageTime = t
	return nil
}
