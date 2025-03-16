package repository_impl

import (
	"context"
	"database/sql"
	"github.com/toormi/go-clean-arch/internal/article-domain/domain/entity"
	"github.com/toormi/go-clean-arch/internal/article-persistence/convert"
	"github.com/toormi/go-clean-arch/internal/article-persistence/dao"
)

type AuthorRepositoryImpl struct {
	*dao.AuthorDAO
	*convert.AuthorConvert
}

// NewAuthorRepositoryImpl will create an implementation of author.Repository
func NewAuthorRepositoryImpl(db *sql.DB) *AuthorRepositoryImpl {
	return &AuthorRepositoryImpl{
		AuthorDAO:     dao.NewAuthorDAO(db),
		AuthorConvert: convert.NewAuthorConvert(),
	}
}

func (m *AuthorRepositoryImpl) GetByID(ctx context.Context, id int64) (entity.Author, error) {
	daoRes, err := m.AuthorDAO.GetByID(ctx, id)
	if err != nil {
		return entity.Author{}, err
	}
	res := *m.AuthorConvert.FromData(&daoRes)
	return res, nil
}
