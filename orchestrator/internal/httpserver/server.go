package httpserver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/tyagivivek8/krc-system/orchestrator/internal/clients"
	"github.com/tyagivivek8/krc-system/orchestrator/internal/llm"
	"github.com/tyagivivek8/krc-system/orchestrator/internal/types"
)

type Server struct {
	queryClient *clients.QueryClient
	auditClient *clients.AuditClient
	llmModel    string
}

func NewServer(qc *clients.QueryClient, ac *clients.AuditClient, llmModel string) *Server {
	return &Server{
		queryClient: qc,
		auditClient: ac,
		llmModel:    llmModel,
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/query", s.handleQuery)
	return mux
}

func (s *Server) handleQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method now allowed", http.StatusMethodNotAllowed)
		return
	}

	var req types.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.Query == "" {
		http.Error(w, "query is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	searchRes, err := s.queryClient.Search(ctx, req.UserId, req.Query, 5)
	if err != nil {
		log.Printf("query search error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if len(searchRes.Contents) == 0 {
		resp := types.QueryResponse{
			Answer:    "I could not find any relevant information in the indexed documents.",
			Citations: nil,
		}
		s.writeJson(w, resp)
		go s.logAudit(req.UserId, req.Query, "", "NO_CONTEXT")
		return
	}

	prompt := llm.BuildPrompt(req.Query, searchRes.Citations, searchRes.Contents)

	answer, err := llm.GenerateAnswer(ctx, s.llmModel, prompt)
	if err != nil {
		log.Printf("llm error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		go s.logAudit(req.UserId, req.Query, "", "LLM_ERROR")
		return
	}

	resp := types.QueryResponse{
		Answer:    answer,
		Citations: searchRes.Citations,
	}
	s.writeJson(w, resp)

	go s.logAudit(req.UserId, req.Query, "", "SUCCESS")
}

func (s *Server) writeJson(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(payload); err != nil {
		log.Printf("writeJSON error: %v", err)
	}
}

func (s *Server) logAudit(userId, query, responseRef, status string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.auditClient.Log(ctx, userId, query, responseRef, status); err != nil {
		log.Printf("audit log error: %v", err)
	}
}
