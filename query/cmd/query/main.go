package main

import (
	"log"
	"net"
	"os"

	"github.com/tyagivivek8/krc-system/query/internal/service"
	querypb "github.com/tyagivivek8/krc-system/query/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	addr := getEnv("QUERY_GRPC_ADDR", ":50052")

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}

	srv := grpc.NewServer()
	querySvc := service.NewService()

	querypb.RegisterQueryServiceServer(srv, querySvc)
	reflection.Register(srv)

	log.Printf("Query gRPC server listening on %s", addr)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
