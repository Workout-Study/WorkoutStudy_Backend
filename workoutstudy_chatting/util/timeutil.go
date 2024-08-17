package util

import (
	"log"
	"strings"
	"time"
)

// ParseMessageTime 함수는 시간 문자열을 입력받아 time.Time 객체를 반환합니다.
func ParseMessageTime(timeStr string) (time.Time, error) {
	// 시간 문자열에 + 또는 - 기호가 없을 때 처리
	if strings.Contains(timeStr, " ") && !strings.Contains(timeStr, "+") && !strings.Contains(timeStr, "-") {
		// 공백 이후 부분을 타임존으로 인식하고 + 기호를 추가
		timeStr = strings.TrimSpace(timeStr[:len(timeStr)-6]) + "+" + timeStr[len(timeStr)-5:]
	}

	// 커스텀 레이아웃
	const customLayout = "2006-01-02 15:04:05.999999-07:00"

	// 커스텀 레이아웃을 사용하여 시간 문자열 파싱
	parsedTime, err := time.Parse(customLayout, timeStr)
	if err != nil {
		log.Printf("시간 파싱 실패: %v", err)
		return time.Time{}, err
	}

	log.Printf("파싱된 시간: %s", parsedTime)
	return parsedTime, nil
}
