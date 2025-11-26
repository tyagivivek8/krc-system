// query/internal/embeddings/ollama.go
package embeddings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

	url := "http://host.docker.internal:11434/api/embeddings"

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("embedding request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	log.Printf("Ollama embeddings response status=%d body=%s", resp.StatusCode, string(raw))

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ollama status: %s", resp.Status)
	}

	var parsed ollamaResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("decoding ollama response: %w", err)
	}

	if len(parsed.Embedding) == 0 {
		return nil, fmt.Errorf("ollama returned empty embedding")
	}

	return parsed.Embedding, nil
}
