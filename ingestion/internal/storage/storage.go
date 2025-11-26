package storage

import (
	"database/sql"
	"embed"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaFS embed.FS

func InitDB(path string) (*sql.DB, error) {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(string(schema)); err != nil {
		return nil, err
	}

	return db, nil
}

func InsertChunk(db *sql.DB, docName string, index int, content string) error {
	_, err := db.Exec(`
        INSERT INTO chunks (document_name, chunk_index, content)
        VALUES (?, ?, ?)
    `, docName, index, content)

	return err
}
