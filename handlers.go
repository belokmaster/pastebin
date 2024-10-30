package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// createPasteHandler creates a new paste and returns a link
func createPasteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Content cannot be empty", http.StatusBadRequest)
		return
	}

	p := Paste{
		ID:        generateID(),
		Content:   content,
		CreatedAt: time.Now(),
	}

	_, err := db.Exec("INSERT INTO pastes (id, content, created_at) VALUES (?, ?, ?)", p.ID, p.Content, p.CreatedAt)
	if err != nil {
		log.Printf("Error inserting paste: %v", err)
		http.Error(w, "Could not save paste. Please try again later.", http.StatusInternalServerError)
		return
	}

	host := r.Host
	if host == "" {
		host = "localhost:6386"
	}
	link := fmt.Sprintf("http://%s/get?id=%s", host, p.ID)

	tmpl, err := template.ParseFiles("templates/paste_created.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Could not load template. Please try again later.", http.StatusInternalServerError)
		return
	}

	data := struct {
		Link string
	}{
		Link: link,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Could not execute template. Please try again later.", http.StatusInternalServerError)
	}
}

// getPasteHandler retrieves a paste by ID and displays it
func getPasteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	var p Paste
	err := db.QueryRow("SELECT id, content, created_at FROM pastes WHERE id = ?", id).Scan(&p.ID, &p.Content, &p.CreatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "Paste not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error querying paste: %v", err)
		http.Error(w, "Could not retrieve paste. Please try again later.", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/view_paste.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Could not load template. Please try again later.", http.StatusInternalServerError)
		return
	}

	data := struct {
		ID        string
		Content   string
		CreatedAt string
	}{
		ID:        p.ID,
		Content:   p.Content,
		CreatedAt: p.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Could not execute template. Please try again later.", http.StatusInternalServerError)
	}
}
