package main

import (
	"database/sql"
	"log"
	"time"
)

// Paste struct for storing text and creation time
type Paste struct {
	ID           string    `json:"id"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	ExpirationAt time.Time `json:"expiration_at"`
}

// savePaste saves a paste to the database
func savePaste(db *sql.DB, p Paste) error {
	query := `INSERT INTO pastes (id, content, created_at) VALUES (?, ?, ?)`
	_, err := db.Exec(query, p.ID, p.Content, p.CreatedAt, p.ExpirationAt)
	return err
}

// getPasteByID retrieves a paste by ID from the database
func getPasteByID(db *sql.DB, id string) (Paste, error) {
	query := `SELECT id, content, created_at FROM pastes WHERE id = ?`
	row := db.QueryRow(query, id)

	var p Paste
	err := row.Scan(&p.ID, &p.Content, &p.CreatedAt)
	return p, err
}

func deleteExpiredPastes(db *sql.DB) error {
	query := `DELETE FROM pastes WHERE expiration_at < ?`
	_, err := db.Exec(query, time.Now())
	return err
}

func startExpirationCleaner(db *sql.DB) {
	ticker := time.NewTicker(1 * time.Minute) // Adjust the interval as needed
	defer ticker.Stop()

	for {
		<-ticker.C
		if err := deleteExpiredPastes(db); err != nil {
			log.Printf("Error deleting expired pastes: %v", err)
		}
	}
}
