package util

import (
	"log"
	"time"
)

// ParseMessageTime 함수는 시간 문자열을 입력받아 time.Time 객체를 반환합니다.
func ParseMessageTime(timeStr string) (time.Time, error) {
	// PostgreSQL의 timestamp with time zone 형식에 맞춘 커스텀 레이아웃
	const customLayout = "2006-01-02 15:04:05.000"

	// 커스텀 레이아웃을 사용하여 시간 문자열 파싱
	parsedTime, err := time.Parse(customLayout, timeStr)
	if err != nil {
		log.Printf("시간 파싱 실패: %v", err)
		return time.Time{}, err
	}

	log.Printf("파싱된 시간: %s", parsedTime)

	return parsedTime, nil
}
