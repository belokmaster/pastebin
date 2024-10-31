package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

// showFormHandler shows the HTML form for inputting text
func showFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("tmpl/create_paste.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Could not load template. Please try again later.", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Could not execute template. Please try again later.", http.StatusInternalServerError)
	}
}

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

	expirationStr := r.FormValue("expiration")
	expiration, err := strconv.Atoi(expirationStr)
	if err != nil {
		http.Error(w, "Invalid expiration value", http.StatusBadRequest)
		return
	}

	expirationAt := time.Now().Add(time.Duration(expiration) * time.Minute)

	p := Paste{
		ID:           generateID(),
		Content:      content,
		CreatedAt:    time.Now(),
		ExpirationAt: expirationAt,
	}

	if err := savePaste(db, p); err != nil {
		log.Printf("Error saving paste: %v", err)
		http.Error(w, "Could not save paste. Please try again later.", http.StatusInternalServerError)
		return
	}

	host := r.Host
	if host == "" {
		host = "localhost:6386"
	}
	link := fmt.Sprintf("http://%s/get?id=%s", host, p.ID)

	tmpl, err := template.ParseFiles("tmpl/paste_created.html")
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

	paste, err := getPasteByID(db, id)
	if err == sql.ErrNoRows {
		http.Error(w, "Paste not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error retrieving paste: %v", err)
		http.Error(w, "Could not retrieve paste. Please try again later.", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("tmpl/view_paste.html")
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
		ID:        paste.ID,
		Content:   paste.Content,
		CreatedAt: paste.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Could not execute template. Please try again later.", http.StatusInternalServerError)
	}
}
