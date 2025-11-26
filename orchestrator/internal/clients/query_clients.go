package clients

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/tyagivivek8/krc-system/orchestrator/internal/types"
	querypb "github.com/tyagivivek8/krc-system/query/proto"
)

type QueryClient struct {
	client querypb.QueryServiceClient
}

func NreQueryClient(addr string) (*QueryClient, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)

	if err != nil {
		return nil, fmt.Errorf("dial query service: %w", err)
	}

	return &QueryClient{
		client: querypb.NewQueryServiceClient(conn),
	}, nil
}

type SearchResult struct {
	Citations []types.Citation
	Contents  []string
}

func (qc *QueryClient) Search(ctx context.Context, userID, q string, limit int32) (*SearchResult, error) {
	req := &querypb.QueryRequest{
		UserId: userID,
		Query:  q,
	}
	resp, err := qc.client.Search(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("query search: %w", err)
	}

	if len(resp.Hits) == 0 {
		return &SearchResult{}, nil
	}

	hits := resp.Hits
	if int(limit) < len(hits) {
		hits = hits[:limit]
	}

	res := &SearchResult{
		Citations: make([]types.Citation, 0, len(hits)),
		Contents:  make([]string, 0, len(hits)),
	}

	for _, h := range hits {
		res.Citations = append(res.Citations, types.Citation{
			DocumentName: h.GetDocumentName(),
			ChunkIndex:   h.GetChunkIndex(),
			Score:        h.GetScore(),
		})
		res.Contents = append(res.Contents, h.GetContent())
	}
	return res, nil
}
