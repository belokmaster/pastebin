package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
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

func main() {
	rand.Seed(time.Now().UnixNano())
	port := rand.Intn(9000) + 1000

	http.HandleFunc("/", showFormHandler)          // Shows the form for creating a paste
	http.HandleFunc("/create", createPasteHandler) // Creates a paste
	http.HandleFunc("/get", getPasteHandler)       // Gets a paste by ID

	log.Printf("Server running on http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// showFormHandler shows the HTML form for inputting text
func showFormHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the HTML template
	tmpl, err := template.ParseFiles("tmpl/create_paste.html")
	if err != nil {
		log.Printf("Error loading template: %v", err) // Log the detailed error
		http.Error(w, "Could not load template. Please try again later.", http.StatusInternalServerError)
		return
	}

	// Execute the template and pass it to the response writer
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Error executing template: %v", err) // Log the detailed error
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

	// Create a new paste
	p := Paste{
		ID:        generateID(),
		Content:   content,
		CreatedAt: time.Now(),
	}

	mu.Lock()
	pasteStore[p.ID] = p
	mu.Unlock()

	// Generate a link to the paste
	host := r.Host
	if host == "" {
		host = "localhost:6386" // Specify your port if known
	}
	link := fmt.Sprintf("http://%s/get?id=%s", host, p.ID)

	// Parse the HTML template
	tmpl, err := template.ParseFiles("tmpl/paste_created.html")
	if err != nil {
		log.Printf("Error loading template: %v", err) // Log the detailed error
		http.Error(w, "Could not load template. Please try again later.", http.StatusInternalServerError)
		return
	}

	// Execute the template and pass the link to it
	data := struct {
		Link string
	}{
		Link: link,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err) // Log the detailed error
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

	mu.Lock()
	paste, exists := pasteStore[id]
	mu.Unlock()

	if !exists {
		http.Error(w, "Paste not found", http.StatusNotFound)
		return
	}

	// Parse the HTML template
	tmpl, err := template.ParseFiles("tmpl/view_paste.html")
	if err != nil {
		log.Printf("Error loading template: %v", err) // Log the detailed error
		http.Error(w, "Could not load template. Please try again later.", http.StatusInternalServerError)
		return
	}

	// Execute the template and pass the paste data to it
	data := struct {
		ID        string
		Content   string
		CreatedAt string
	}{
		ID:        paste.ID,
		Content:   paste.Content,
		CreatedAt: paste.CreatedAt.Format("2006-01-02 15:04:05"), // Форматирование даты и времени
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err) // Log the detailed error
		http.Error(w, "Could not execute template. Please try again later.", http.StatusInternalServerError)
	}
}

// generateID generates a random ID for the paste
func generateID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := make([]byte, 8)
	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}
	return string(id)
}
