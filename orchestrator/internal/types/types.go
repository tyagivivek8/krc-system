package types

type QueryRequest struct {
	UserId string `json:"user_id"`
	Query  string `json:"query"`
}

type Citation struct {
	DocumentName string  `json:"document_name"`
	ChunkIndex   int32   `json:"chunk_index"`
	Score        float32 `json:"score"`
}

type QueryResponse struct {
	Answer    string     `json:"answer"`
	Citations []Citation `json:"citations"`
}
