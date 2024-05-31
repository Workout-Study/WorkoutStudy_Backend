package model

import "time"

type User struct {
	ID        int
	Nickname  string
	State     bool
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}
