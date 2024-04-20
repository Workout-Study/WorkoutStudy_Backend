// /model/fit_mate.go
package model

import "time"

type FitMate struct {
	ID         int
	FitGroupID int
	Username   string
	Nickname   string
	State      bool
	CreatedAt  time.Time
	CreatedBy  string
	UpdatedAt  time.Time
	UpdatedBy  string
}
