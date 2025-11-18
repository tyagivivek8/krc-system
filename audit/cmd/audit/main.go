package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/tyagivivek8/krc-system/audit/internal/audit"
	"github.com/tyagivivek8/krc-system/audit/internal/db"
	auditpb "github.com/tyagivivek8/krc-system/audit/proto"
)

func main() {
	addr := getEnv("AUDIT_GRPC_ADDR", ":50051")
	dsn := getEnv("AUDIT_DB_DSN", "postgres://audit:auditpass@postgres:5432/auditlog?sslmode=disable")

	database, err := db.New(dsn)
	if err != nil {
		log.Fatalf("Failed to create DB Pool: %v", err)
	}
	defer database.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := database.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping DB %v", err)
	}

	repo := audit.NewRepository(database)
	service := audit.NewService(repo)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}

	grpcServer := grpc.NewServer()

	auditpb.RegisterAuditServiceServer(grpcServer, service)

	//Health checks
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)

	log.Printf("Audit gRPC Server listening on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
