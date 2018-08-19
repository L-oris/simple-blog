package models

import "time"

type Post struct {
	Title     string
	Content   string
	CreatedAt time.Time
}
