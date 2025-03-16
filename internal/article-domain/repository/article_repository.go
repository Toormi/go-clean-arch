package repository

import (
	"context"
	"github.com/toormi/go-clean-arch/internal/article-persistence/persistence"
)

// ArticleRepository represent the article's repository contract
//
//go:generate mockery --name ArticleRepository
type ArticleRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []persistence.Article, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (persistence.Article, error)
	GetByTitle(ctx context.Context, title string) (persistence.Article, error)
	Update(ctx context.Context, ar *persistence.Article) error
	Store(ctx context.Context, a *persistence.Article) error
	Delete(ctx context.Context, id int64) error
}
