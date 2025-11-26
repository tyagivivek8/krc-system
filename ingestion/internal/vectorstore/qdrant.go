package vectorstore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type point struct {
	ID      string                 `json:"id"`
	Vector  []float32              `json:"vector"`
	Payload map[string]interface{} `json:"payload"`
}

type upsertRequest struct {
	Points []point `json:"points"`
}

func UpsertChunk(
	collection string,
	id string,
	vector []float32,
	docName string,
	chunkIndex int,
	content string,
) error {
	baseURL := os.Getenv("QDRANT_URL")
	if baseURL == "" {
		baseURL = "http://qdrant:6333" // inside Docker network
	}

	// Build payload
	payload := map[string]interface{}{
		"document_name": docName,
		"chunk_index":   chunkIndex,
		"content":       content,
	}

	reqBody := upsertRequest{
		Points: []point{
			{
				ID:      id,
				Vector:  vector,
				Payload: payload,
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/collections/%s/points?wait=true", baseURL, collection)

	req, err := http.NewRequest("PUT", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("qdrant upsert status: %s, body: %s", resp.Status, string(b))
	}

	return nil
}
