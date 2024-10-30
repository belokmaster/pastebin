package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
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
