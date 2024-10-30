package main

import (
	"time"
)

// Paste struct for storing text and creation time
type Paste struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
