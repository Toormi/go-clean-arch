package convert

import (
	"github.com/jinzhu/copier"
	"github.com/toormi/go-clean-arch/internal/article-domain/domain/entity"
	"github.com/toormi/go-clean-arch/internal/article-persistence/persistence"
)

type ArticleConvert struct {
}

func NewArticleConvert() *ArticleConvert {
	return &ArticleConvert{}
}

func (a *ArticleConvert) FromData(article *persistence.Article) *entity.Article {
	result := entity.Article{}
	_ = copier.Copy(&result, &article)
	return &result
}

func (a *ArticleConvert) ToData(article *entity.Article) *persistence.Article {
	result := persistence.Article{}
	_ = copier.Copy(&result, &article)
	return &result
}
