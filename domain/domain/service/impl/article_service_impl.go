package impl

import (
	"context"
	"github.com/toormi/go-clean-arch/domain/repository"
	"github.com/toormi/go-clean-arch/persistence/persistence"
	"github.com/toormi/go-clean-arch/web"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	articleRepo repository.ArticleRepository
	authorRepo  repository.AuthorRepository
}

// NewService will create a new article service object
func NewService(a repository.ArticleRepository, ar repository.AuthorRepository) *Service {
	return &Service{
		articleRepo: a,
		authorRepo:  ar,
	}
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
func (a *Service) fillAuthorDetails(ctx context.Context, data []persistence.Article) ([]persistence.Article, error) {
	g, ctx := errgroup.WithContext(ctx)
	// Get the author's id
	mapAuthors := map[int64]persistence.Author{}

	for _, article := range data { //nolint
		mapAuthors[article.Author.ID] = persistence.Author{}
	}
	// Using goroutine to fetch the author's detail
	chanAuthor := make(chan persistence.Author)
	for authorID := range mapAuthors {
		authorID := authorID
		g.Go(func() error {
			res, err := a.authorRepo.GetByID(ctx, authorID)
			if err != nil {
				return err
			}
			chanAuthor <- res
			return nil
		})
	}

	go func() {
		defer close(chanAuthor)
		err := g.Wait()
		if err != nil {
			logrus.Error(err)
			return
		}

	}()

	for author := range chanAuthor {
		if author != (persistence.Author{}) {
			mapAuthors[author.ID] = author
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// merge the author's data
	for index, item := range data { //nolint
		if a, ok := mapAuthors[item.Author.ID]; ok {
			data[index].Author = a
		}
	}
	return data, nil
}

func (a *Service) Fetch(ctx context.Context, cursor string, num int64) (res []persistence.Article, nextCursor string, err error) {
	res, nextCursor, err = a.articleRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	res, err = a.fillAuthorDetails(ctx, res)
	if err != nil {
		nextCursor = ""
	}
	return
}

func (a *Service) GetByID(ctx context.Context, id int64) (res persistence.Article, err error) {
	res, err = a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return persistence.Article{}, err
	}
	res.Author = resAuthor
	return
}

func (a *Service) Update(ctx context.Context, ar *persistence.Article) (err error) {
	ar.UpdatedAt = time.Now()
	return a.articleRepo.Update(ctx, ar)
}

func (a *Service) GetByTitle(ctx context.Context, title string) (res persistence.Article, err error) {
	res, err = a.articleRepo.GetByTitle(ctx, title)
	if err != nil {
		return
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return persistence.Article{}, err
	}

	res.Author = resAuthor
	return
}

func (a *Service) Store(ctx context.Context, m *persistence.Article) (err error) {
	existedArticle, _ := a.GetByTitle(ctx, m.Title) // ignore if any error
	if existedArticle != (persistence.Article{}) {
		return web.ErrConflict
	}

	err = a.articleRepo.Store(ctx, m)
	return
}

func (a *Service) Delete(ctx context.Context, id int64) (err error) {
	existedArticle, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedArticle == (persistence.Article{}) {
		return web.ErrNotFound
	}
	return a.articleRepo.Delete(ctx, id)
}
