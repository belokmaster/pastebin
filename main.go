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

	var err error
	db, err = sql.Open("sqlite3", "./pastebin.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = createTable()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", showFormHandler)
	http.HandleFunc("/create", createPasteHandler)
	http.HandleFunc("/get", getPasteHandler)

	log.Printf("Server running on http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS pastes (
		id TEXT PRIMARY KEY,
		content TEXT,
		created_at DATETIME
	);
	`
	_, err := db.Exec(query)
	return err
}
