package main

import (
	"log"
	"net/http"

	"github.com/tyagivivek8/krc-system/ingestion/internal/handlers"
	"github.com/tyagivivek8/krc-system/ingestion/internal/storage"
)

func main() {
	db, err := storage.InitDB("data/chunks.db")
	if err != nil {
		log.Fatalf("DB Init error: %v", err)
	}
	handlers.SetDB(db)

	mux := http.NewServeMux()

	mux.HandleFunc("/upload", handlers.UploadHandler)

	log.Println("Ingestion running on :8001")
	if err := http.ListenAndServe(":8001", mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
