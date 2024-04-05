// util/timeutil.go
package util

import (
	"time"
)

// ConvertUTCToLocalTime 함수가 time.Time 객체를 반환하도록 수정
func ConvertUTCToLocalTime(utcTimeStr, locationStr string) (time.Time, error) {
	utcTime, err := time.Parse(time.RFC3339, utcTimeStr)
	if err != nil {
		return time.Time{}, err
	}

	loc, err := time.LoadLocation(locationStr)
	if err != nil {
		return time.Time{}, err
	}

	localTime := utcTime.In(loc)
	return localTime, nil
}
