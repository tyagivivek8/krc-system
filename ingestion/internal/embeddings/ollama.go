package embeddings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ollamaResponse struct {
	Embedding []float32 `json:"embedding"`
}

func GenerateEmbedding(text string) ([]float32, error) {
	body, _ := json.Marshal(ollamaRequest{
		Model:  "nomic-embed-text",
		Prompt: text,
	})

	req, err := http.NewRequest("POST", "http://host.docker.internal:11434/api/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("embedding request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ollama status: %s", resp.Status)
	}

	var parsed ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	return parsed.Embedding, nil
}
