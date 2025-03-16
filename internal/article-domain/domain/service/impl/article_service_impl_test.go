package impl_test

import (
	"context"
	"errors"
	"github.com/toormi/go-clean-arch/internal/article-domain/domain/service/impl"
	mocks2 "github.com/toormi/go-clean-arch/internal/article-domain/repository/mocks"
	persistence2 "github.com/toormi/go-clean-arch/internal/article-persistence/persistence"
	"github.com/toormi/go-clean-arch/internal/server"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetchArticle(t *testing.T) {
	mockArticleRepo := new(mocks2.ArticleRepository)
	mockArticle := persistence2.Article{
		Title:   "Hello",
		Content: "Content",
	}

	mockListArtilce := make([]persistence2.Article, 0)
	mockListArtilce = append(mockListArtilce, mockArticle)

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(mockListArtilce, "next-cursor", nil).Once()
		mockAuthor := persistence2.Author{
			ID:   1,
			Name: "Iman Tumorang",
		}
		mockAuthorrepo := new(mocks2.AuthorRepository)
		mockAuthorrepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)
		u := impl.NewService(mockArticleRepo, mockAuthorrepo)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)
		cursorExpected := "next-cursor"
		assert.Equal(t, cursorExpected, nextCursor)
		assert.NotEmpty(t, nextCursor)
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListArtilce))

		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(nil, "", errors.New("Unexpexted Error")).Once()

		mockAuthorrepo := new(mocks2.AuthorRepository)
		u := impl.NewService(mockArticleRepo, mockAuthorrepo)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)

		assert.Empty(t, nextCursor)
		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	mockArticleRepo := new(mocks2.ArticleRepository)
	mockArticle := persistence2.Article{
		Title:   "Hello",
		Content: "Content",
	}
	mockAuthor := persistence2.Author{
		ID:   1,
		Name: "Iman Tumorang",
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockArticle, nil).Once()
		mockAuthorrepo := new(mocks2.AuthorRepository)
		mockAuthorrepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)
		u := impl.NewService(mockArticleRepo, mockAuthorrepo)

		a, err := u.GetByID(context.TODO(), mockArticle.ID)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(persistence2.Article{}, errors.New("Unexpected")).Once()

		mockAuthorrepo := new(mocks2.AuthorRepository)
		u := impl.NewService(mockArticleRepo, mockAuthorrepo)

		a, err := u.GetByID(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		assert.Equal(t, persistence2.Article{}, a)

		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
}

func TestStore(t *testing.T) {
	mockArticleRepo := new(mocks2.ArticleRepository)
	mockArticle := persistence2.Article{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		tempMockArticle := mockArticle
		tempMockArticle.ID = 0
		mockArticleRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(persistence2.Article{}, server.ErrNotFound).Once()
		mockArticleRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.Article")).Return(nil).Once()

		mockAuthorrepo := new(mocks2.AuthorRepository)
		u := impl.NewService(mockArticleRepo, mockAuthorrepo)

		err := u.Store(context.TODO(), &tempMockArticle)

		assert.NoError(t, err)
		assert.Equal(t, mockArticle.Title, tempMockArticle.Title)
		mockArticleRepo.AssertExpectations(t)
	})
	t.Run("existing-title", func(t *testing.T) {
		existingArticle := mockArticle
		mockArticleRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(existingArticle, nil).Once()
		mockAuthor := persistence2.Author{
			ID:   1,
			Name: "Iman Tumorang",
		}
		mockAuthorrepo := new(mocks2.AuthorRepository)
		mockAuthorrepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)

		u := impl.NewService(mockArticleRepo, mockAuthorrepo)

		err := u.Store(context.TODO(), &mockArticle)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockArticleRepo := new(mocks2.ArticleRepository)
	mockArticle := persistence2.Article{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockArticle, nil).Once()

		mockArticleRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		mockAuthorrepo := new(mocks2.AuthorRepository)
		u := impl.NewService(mockArticleRepo, mockAuthorrepo)

		err := u.Delete(context.TODO(), mockArticle.ID)

		assert.NoError(t, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("article-is-not-exist", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(persistence2.Article{}, nil).Once()

		mockAuthorrepo := new(mocks2.AuthorRepository)
		u := impl.NewService(mockArticleRepo, mockAuthorrepo)

		err := u.Delete(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("error-happens-in-db", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(persistence2.Article{}, errors.New("Unexpected Error")).Once()

		mockAuthorrepo := new(mocks2.AuthorRepository)
		u := impl.NewService(mockArticleRepo, mockAuthorrepo)

		err := u.Delete(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockArticleRepo := new(mocks2.ArticleRepository)
	mockArticle := persistence2.Article{
		Title:   "Hello",
		Content: "Content",
		ID:      23,
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("Update", mock.Anything, &mockArticle).Once().Return(nil)

		mockAuthorrepo := new(mocks2.AuthorRepository)
		u := impl.NewService(mockArticleRepo, mockAuthorrepo)

		err := u.Update(context.TODO(), &mockArticle)
		assert.NoError(t, err)
		mockArticleRepo.AssertExpectations(t)
	})
}
