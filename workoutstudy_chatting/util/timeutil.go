package util

import (
	"log"
	"time"
)

// ParseMessageTime 함수는 시간 문자열을 입력받아 time.Time 객체를 반환합니다.
func ParseMessageTime(timeStr string) (time.Time, error) {
	// 입력된 timeStr 로그로 출력
	log.Printf("Original timeStr in ParseMessageTime: %s", timeStr)

	// 커스텀 레이아웃 정의: 타임존 정보가 있는 경우를 처리합니다.
	const customLayout = "2006-01-02 15:04:05.999999-07:00"

	// 커스텀 레이아웃을 사용하여 시간 문자열을 파싱합니다.
	parsedTime, err := time.Parse(customLayout, timeStr)
	if err != nil {
		log.Printf("시간 파싱 실패: %v", err)
		return time.Time{}, err
	}

	// 파싱된 시간 로그로 출력
	log.Printf("Parsed time in ParseMessageTime: %s", parsedTime)

	return parsedTime, nil
}
