package repository

import (
	"context"
	"github.com/toormi/go-clean-arch/internal/article-domain/domain/entity"
)

// AuthorRepository represent the author's repository contract
//
//go:generate mockery --name AuthorRepository
type AuthorRepository interface {
	GetByID(ctx context.Context, id int64) (entity.Author, error)
}
