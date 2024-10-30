package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	rand.Seed(time.Now().UnixNano())
	port := rand.Intn(9000) + 1000

	// Initialize the database
	var err error
	db, err = sql.Open("sqlite3", "./pastes.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := initializeDatabase(db); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	http.HandleFunc("/", showFormHandler)          // Shows the form for creating a paste
	http.HandleFunc("/create", createPasteHandler) // Creates a paste
	http.HandleFunc("/get", getPasteHandler)       // Gets a paste by ID

	log.Printf("Server running on http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func initializeDatabase(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS pastes (
		id TEXT PRIMARY KEY,
		content TEXT,
		created_at TIMESTAMP
	);`
	_, err := db.Exec(query)
	return err
}
