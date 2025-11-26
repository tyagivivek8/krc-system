package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tyagivivek8/krc-system/orchestrator/internal/clients"
	"github.com/tyagivivek8/krc-system/orchestrator/internal/httpserver"
)

func main() {
	httpAddr := getEnv("ORCH_HTTP_ADDR", ":8080")
	queryAddr := getEnv("QUERY_GRPC_ADDR", "query:50052")
	auditAddr := getEnv("AUDIT_GRPC_ADDR", "audit:50051")
	llmModel := getEnv("LLM_MODEL", "llama3")

	log.Printf("Connecting to Query service at %s", queryAddr)
	qc, err := clients.NreQueryClient(queryAddr)
	if err != nil {
		log.Fatalf("init query client: %v", err)
	}

	log.Printf("Connecting to audit service at %s", auditAddr)
	ac, err := clients.NewAuditClient(auditAddr)
	if err != nil {
		log.Fatalf("init audit client: %v", err)
	}

	srv := httpserver.NewServer(qc, ac, llmModel)

	log.Printf("Orchestrator HTTP server listening on %s", httpAddr)
	if err := http.ListenAndServe(httpAddr, srv.Routes()); err != nil {
		log.Fatalf("http server failed: %v", err)
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
