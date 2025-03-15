package service

import (
	"context"
	"github.com/toormi/go-clean-arch/persistence/persistence"
)

// ArticleService represent the article's usecases
//
//go:generate mockery --name ArticleService
type ArticleService interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]persistence.Article, string, error)
	GetByID(ctx context.Context, id int64) (persistence.Article, error)
	Update(ctx context.Context, ar *persistence.Article) error
	GetByTitle(ctx context.Context, title string) (persistence.Article, error)
	Store(context.Context, *persistence.Article) error
	Delete(ctx context.Context, id int64) error
}
