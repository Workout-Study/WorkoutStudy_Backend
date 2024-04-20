// /model/fit_group.go
package model

import "time"

type FitGroup struct {
	ID           int
	FitGroupName string
	MaxFitMate   int
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time
	UpdatedBy    string
}
