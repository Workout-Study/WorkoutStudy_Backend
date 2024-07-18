package model

import "time"

type FitMate struct {
	ID         int
	UserID     int
	FitGroupID int
	State      bool
	CreatedAt  time.Time
	CreatedBy  string
	UpdatedAt  time.Time
	UpdatedBy  string
}
