package clients

import (
	"context"
	"fmt"
	"time"

	auditpb "github.com/tyagivivek8/krc-system/audit/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuditClient struct {
	client auditpb.AuditServiceClient
}

func NewAuditClient(addr string) (*AuditClient, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("dial audit service: %w", err)
	}
	return &AuditClient{
		client: auditpb.NewAuditServiceClient(conn),
	}, nil
}

func (ac *AuditClient) Log(ctx context.Context, userID, query, responseRef, status string) error {
	req := &auditpb.LogQueryRequest{
		UserId:      userID,
		Query:       query,
		ResponseRef: responseRef,
		Status:      status,
	}

	_, err := ac.client.LogQuery(ctx, req)
	if err != nil {
		return fmt.Errorf("audit log query: %w", err)
	}
	return nil
}
