package main

import (
	"sync"
	"time"
)

// Paste struct for storing text and creation time
type Paste struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// In-memory storage for pastes
var (
	pasteStore = make(map[string]Paste)
	mu         sync.Mutex
)
