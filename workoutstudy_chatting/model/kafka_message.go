package model

import "time"

type UserCreateEvent struct {
	UserID    int       `json:"userId"`
	Nickname  string    `json:"nickname"`
	State     bool      `json:"state"`
	ImageUrl  string    `json:"imageUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// 유저 생성 이벤트
// user-create-event
// {
// "id": "user_id",
// "nickname": "nickname",
// "state": false
// "createdAt": 시간데이터를 문자열로
// "updatedAt": 시간데이터를 문자열로
// }

// user-info-event
// zero payload 용 이벤트 토픽
// kafka message 의 value 에 key:value 형태가 아닌 단순 숫자타입 데이터(user_id) 하나만 전송
