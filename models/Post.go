package models

import (
	"time"
)

// Post model
type Post struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

// IsValidPost validates Post coming from client
func IsValidPost(p Post) bool {
	if p.ID == "" || p.Title == "" {
		return false
	}
	return true
}
