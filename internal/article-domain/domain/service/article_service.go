package service

import (
	"context"
	"github.com/toormi/go-clean-arch/internal/article-domain/domain/entity"
)

// ArticleService represent the article's usecases
//
//go:generate mockery --name ArticleService
type ArticleService interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]entity.Article, string, error)
	GetByID(ctx context.Context, id int64) (entity.Article, error)
	Update(ctx context.Context, ar *entity.Article) error
	GetByTitle(ctx context.Context, title string) (entity.Article, error)
	Store(context.Context, *entity.Article) error
	Delete(ctx context.Context, id int64) error
}
