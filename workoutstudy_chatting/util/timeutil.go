package util

import (
	"log"
	"time"
)

// ParseMessageTime 함수는 시간 문자열을 입력받아 time.Time 객체를 반환합니다.
func ParseMessageTime(timeStr string) (time.Time, error) {
	// 커스텀 레이아웃: 타임존이 있는 경우를 처리합니다.
	const customLayout = "2006-01-02 15:04:05.999999-07:00"

	// 커스텀 레이아웃을 사용하여 시간 문자열을 파싱합니다.
	parsedTime, err := time.Parse(customLayout, timeStr)
	if err != nil {
		log.Printf("시간 파싱 실패: %v", err)
		return time.Time{}, err
	}

	log.Printf("파싱된 시간: %s", parsedTime)
	return parsedTime, nil
}
