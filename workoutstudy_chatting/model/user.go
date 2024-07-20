package model

import "time"

type Users struct {
	ID        int
	Nickname  string
	State     bool
	ImageUrl  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
