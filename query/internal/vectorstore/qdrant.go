package vectorstore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type searchRequest struct {
	Vector      []float32 `json:"vector"`
	Limit       int       `json:"limit"`
	WithPayload bool      `json:"with_payload"`
}

type searchResultPoint struct {
	ID      interface{}            `json:"id"`
	Score   float32                `json:"score"`
	Payload map[string]interface{} `json:"payload"`
}

type searchResponse struct {
	Result []searchResultPoint `json:"result"`
	Status string              `json:"status"`
}

type Hit struct {
	ID           string
	DocumentName string
	ChunkIndex   int32
	Content      string
	Score        float32
}

func Search(collection string, vector []float32, limit int) ([]Hit, error) {
	baseURL := os.Getenv("QDRANT_URL")
	if baseURL == "" {
		baseURL = "http://qdrant:6333"
	}

	reqBody := searchRequest{
		Vector:      vector,
		Limit:       limit,
		WithPayload: true,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/collections/%s/points/search", baseURL, collection)

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("qdrant search status: %s", resp.Status)
	}

	var parsed searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	hits := make([]Hit, 0, len(parsed.Result))
	for _, p := range parsed.Result {
		h := Hit{
			Score: p.Score,
		}
		switch v := p.ID.(type) {
		case string:
			h.ID = v
		default:
			h.ID = fmt.Sprint(v)
		}

		if dn, ok := p.Payload["document_name"].(string); ok {
			h.DocumentName = dn
		}
		if ci, ok := p.Payload["chunk_index"].(float64); ok {
			h.ChunkIndex = int32(ci)
		}
		if c, ok := p.Payload["content"].(string); ok {
			h.Content = c
		}

		hits = append(hits, h)
	}
	return hits, nil
}
