package audit

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tyagivivek8/krc-system/audit/internal/db"
)

type Log struct {
	ID          uuid.UUID
	UserID      string
	Query       string
	ResponseRef string
	Status      string
	Timestamp   time.Time
}

type Repository struct {
	db *db.DB
}

func NewRepository(d *db.DB) *Repository {
	return &Repository{db: d}
}

func (r *Repository) InsertLog(ctx context.Context, l *Log) error {
	const q = `
		INSERT INTO audit_logs (id, user_id, query, response_ref, status)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Pool.Exec(
		ctx,
		q,
		l.ID,
		l.UserID,
		l.Query,
		l.ResponseRef,
		l.Status,
	)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Log, error) {
	const q = `
		SELECT id, user_id, query, response_ref, status, timestamp
		FROM audit_logs
		WHERE id = $1
	`
	row := r.db.Pool.QueryRow(ctx, q, id)

	var l Log
	err := row.Scan(
		&l.ID,
		&l.UserID,
		&l.Query,
		&l.ResponseRef,
		&l.Status,
		&l.Timestamp,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &l, nil
}
