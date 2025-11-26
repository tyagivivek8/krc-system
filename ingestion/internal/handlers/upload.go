package handlers

import (
	"database/sql"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/tyagivivek8/krc-system/ingestion/internal/chunker"
	"github.com/tyagivivek8/krc-system/ingestion/internal/embeddings"
	"github.com/tyagivivek8/krc-system/ingestion/internal/extractor"
	"github.com/tyagivivek8/krc-system/ingestion/internal/storage"
	"github.com/tyagivivek8/krc-system/ingestion/internal/vectorstore"
)

var db *sql.DB

func SetDB(d *sql.DB) {
	db = d
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file missing", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Received file: %s\n", header.Filename)

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Read error", http.StatusInternalServerError)
		return
	}

	text, err := extractor.ExtractText(data, header.Filename)
	if err != nil {
		http.Error(w, "Extract error", http.StatusInternalServerError)
		return
	}

	chunks := chunker.Chunk(text, 1024)
	for i, chunk := range chunks {
		err := storage.InsertChunk(db, header.Filename, i, chunk)
		if err != nil {
			http.Error(w, "storage error", http.StatusInternalServerError)
			return
		}

		embedding, err := embeddings.GenerateEmbedding(chunk)
		if err != nil {
			log.Printf("embedding error for chunk %d: %v", i, err)
			continue
		}

		id := uuid.New().String()

		if err := vectorstore.UpsertChunk(
			"krc_docs",
			id,
			embedding,
			header.Filename,
			i,
			chunk,
		); err != nil {
			log.Printf("qdrant upsert error for chunk %d: %v", i, err)
			continue
		}

	}
	log.Printf("Extracted %d chunks\n", len(chunks))

	w.Write([]byte("OK"))
}
