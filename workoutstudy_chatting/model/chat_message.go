package model

import (
	"encoding/json"
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
	FitGroupID  int         `json:"fitGroupId"`
	FitMateID   int         `json:"fitMateId"`
	Message     string      `json:"message"`
	MessageTime time.Time   `json:"messageTime"`
	MessageType MessageType `json:"messageType"`
}

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

	// 따옴표를 제거한 시간 형식 문자열 사용
	t, err := time.Parse("2006-01-02T15:04:05.999999999", tmp.MessageTime)
	if err != nil {
		return err
	}

	cm.MessageTime = t
	return nil
}
