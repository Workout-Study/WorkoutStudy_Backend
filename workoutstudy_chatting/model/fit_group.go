// /model/fit_group.go
package model

import "time"

type FitGroup struct {
	ID                  int
	FitLeaderUserID     int
	FitGroupName        string
	Category            int
	Cycle               int  // 운동 인증 주기 ( 1: 일주일, 2: 한달, 3: 일년 )
	Frequency           int  // 주기별 운동 인증 필요 횟수
	PresentFitMateCount int  // 현재 fit group에 속한 fit mate 수
	MaxFitMate          int  // fit group의 최대 fit mate 수
	State               bool // fit group의 상태 (false: 활성, true: 비활성)
	CreatedAt           time.Time
	CreatedBy           string
	UpdatedAt           time.Time
	UpdatedBy           string
}
