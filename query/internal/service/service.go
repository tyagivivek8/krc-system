package service

import (
	"context"
	"log"

	"github.com/tyagivivek8/krc-system/query/internal/embeddings"
	"github.com/tyagivivek8/krc-system/query/internal/vectorstore"
	querypb "github.com/tyagivivek8/krc-system/query/proto"
)

type Service struct {
	querypb.UnimplementedQueryServiceServer
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Search(ctx context.Context, req *querypb.QueryRequest) (*querypb.SearchResponse, error) {
	vec, err := embeddings.GenerateEmbedding(req.GetQuery())
	if err != nil {
		log.Printf("embedding error: %v", err)
		return nil, err
	}

	hits, err := vectorstore.Search("krc_docs", vec, 5)
	if err != nil {
		log.Printf("qdrant search error: %v", err)
		return nil, err
	}

	resp := &querypb.SearchResponse{}
	for _, h := range hits {
		resp.Hits = append(resp.Hits, &querypb.SearchHit{
			Id:           h.ID,
			DocumentName: h.DocumentName,
			ChunkIndex:   h.ChunkIndex,
			Content:      h.Content,
			Score:        h.Score,
		})
	}
	return resp, nil
}
