package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tyagivivek8/krc-system/orchestrator/internal/types"
)

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

func BuildPrompt(q string, hits []types.Citation, contexts []string) string {
	var sb strings.Builder

	sb.WriteString("You are a helpful assistant answering questions based only on the given context.\n")
	sb.WriteString("If the answer cannot be found in the context, say you are not sure.\n\n")
	sb.WriteString("Context:\n")

	for i, c := range contexts {
		sb.WriteString(fmt.Sprintf("[Chunk %d from %s]\n", hits[i].ChunkIndex, hits[i].DocumentName))
		sb.WriteString(c)
		sb.WriteString("\n\n")
	}

	sb.WriteString("Question:\n")
	sb.WriteString(q)
	sb.WriteString("\nAnswer:\n")

	return sb.String()
}

func GenerateAnswer(ctx context.Context, model string, prompt string) (string, error) {
	body, _ := json.Marshal(ollamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", "http://host.docker.internal:11434/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("ollama status: %s body: %s", resp.Status, string(raw))
	}

	var parsed ollamaResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("decoding ollama response: %w", err)
	}
	return strings.TrimSpace(parsed.Response), nil
}
