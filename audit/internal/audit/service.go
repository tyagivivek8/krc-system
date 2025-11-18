package audit

import (
	"context"
	"time"

	"github.com/google/uuid"
	auditpb "github.com/tyagivivek8/krc-system/audit/proto"
)

type Service struct {
	auditpb.UnimplementedAuditServiceServer
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) LogQuery(ctx context.Context, req *auditpb.LogQueryRequest) (*auditpb.LogQueryResponse, error) {
	// Add timeout on DB call
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	id := uuid.New()

	logEntry := &Log{
		ID:          id,
		UserID:      req.GetUserId(),
		Query:       req.GetQuery(),
		ResponseRef: req.GetResponseRef(),
		Status:      req.GetStatus(),
	}

	if err := s.repo.InsertLog(ctx, logEntry); err != nil {
		return nil, err
	}

	return &auditpb.LogQueryResponse{
		Id:     id.String(),
		Status: logEntry.Status,
	}, nil
}
