package dao

import (
	"context"
	"database/sql"
	"github.com/toormi/go-clean-arch/internal/article-persistence/persistence"
)

type AuthorDAO struct {
	DB *sql.DB
}

// NewAuthorDAO will create an implementation of author.Repository
func NewAuthorDAO(db *sql.DB) *AuthorDAO {
	return &AuthorDAO{
		DB: db,
	}
}

func (m *AuthorDAO) getOne(ctx context.Context, query string, args ...interface{}) (res persistence.Author, err error) {
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return persistence.Author{}, err
	}
	row := stmt.QueryRowContext(ctx, args...)
	res = persistence.Author{}

	err = row.Scan(
		&res.ID,
		&res.Name,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	return
}

func (m *AuthorDAO) GetByID(ctx context.Context, id int64) (persistence.Author, error) {
	query := `SELECT id, name, created_at, updated_at FROM author WHERE id=?`
	return m.getOne(ctx, query, id)
}
