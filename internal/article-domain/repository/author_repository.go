package repository

import (
	"context"
	"github.com/toormi/go-clean-arch/internal/article-persistence/persistence"
)

// AuthorRepository represent the author's repository contract
//
//go:generate mockery --name AuthorRepository
type AuthorRepository interface {
	GetByID(ctx context.Context, id int64) (persistence.Author, error)
}
