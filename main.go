package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", NewHandler()); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
