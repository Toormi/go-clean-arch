package repository_impl

import (
	"context"
	"database/sql"
	"github.com/toormi/go-clean-arch/internal/article-domain/domain/entity"
	"github.com/toormi/go-clean-arch/internal/article-persistence/convert"
	"github.com/toormi/go-clean-arch/internal/article-persistence/dao"
)

type ArticleRepositoryImpl struct {
	*dao.ArticleDAO
	*convert.ArticleConvert
}

// NewArticleRepositoryImpl will create an object that represent the article.Repository interface
func NewArticleRepositoryImpl(conn *sql.DB) *ArticleRepositoryImpl {
	return &ArticleRepositoryImpl{
		ArticleDAO:     dao.NewArticleDAO(conn),
		ArticleConvert: convert.NewArticleConvert(),
	}
}

func (m *ArticleRepositoryImpl) Fetch(ctx context.Context, cursor string, num int64) (res []entity.Article, nextCursor string, err error) {
	daoRes, nextCursor, err := m.ArticleDAO.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}
	res = make([]entity.Article, len(daoRes))
	for i, v := range daoRes {
		res[i] = *m.ArticleConvert.FromData(&v)
	}
	return res, nextCursor, nil
}
func (m *ArticleRepositoryImpl) GetByID(ctx context.Context, id int64) (res entity.Article, err error) {
	daoRes, err := m.ArticleDAO.GetByID(ctx, id)
	if err != nil {
		return entity.Article{}, err
	}
	res = *m.ArticleConvert.FromData(&daoRes)
	return res, nil
}

func (m *ArticleRepositoryImpl) GetByTitle(ctx context.Context, title string) (res entity.Article, err error) {
	daoRes, err := m.ArticleDAO.GetByTitle(ctx, title)
	if err != nil {
		return entity.Article{}, err
	}
	res = *m.ArticleConvert.FromData(&daoRes)
	return res, nil
}

func (m *ArticleRepositoryImpl) Store(ctx context.Context, a *entity.Article) (err error) {
	daoArticle := m.ArticleConvert.ToData(a)
	err = m.ArticleDAO.Store(ctx, daoArticle)
	if err != nil {
		return err
	}
	a.ID = daoArticle.ID
	return nil
}

func (m *ArticleRepositoryImpl) Delete(ctx context.Context, id int64) (err error) {
	return m.ArticleDAO.Delete(ctx, id)
}
func (m *ArticleRepositoryImpl) Update(ctx context.Context, ar *entity.Article) (err error) {
	daoArticle := m.ArticleConvert.ToData(ar)
	err = m.ArticleDAO.Update(ctx, daoArticle)
	if err != nil {
		return err
	}
	ar.ID = daoArticle.ID
	return nil
}
