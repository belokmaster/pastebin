package main

import (
	"math/rand"
)

// generateID generates a random ID for the paste
func generateID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := make([]byte, 8)
	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}
	return string(id)
}
