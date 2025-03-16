package repository

import (
	"context"
	"github.com/bxcodec/go-clean-arch/domain/domain/entity"
)

// ArticleRepository represent the article's repository contract
//
//go:generate mockery --name ArticleRepository
type ArticleRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []entity.Article, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (entity.Article, error)
	GetByTitle(ctx context.Context, title string) (entity.Article, error)
	Update(ctx context.Context, ar *entity.Article) error
	Store(ctx context.Context, a *entity.Article) error
	Delete(ctx context.Context, id int64) error
}
