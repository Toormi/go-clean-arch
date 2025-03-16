package convert

import (
	"github.com/jinzhu/copier"
	"github.com/toormi/go-clean-arch/internal/article-domain/domain/entity"
	"github.com/toormi/go-clean-arch/internal/article-persistence/persistence"
)

type AuthorConvert struct {
}

func NewAuthorConvert() *AuthorConvert {
	return &AuthorConvert{}
}

func (a *AuthorConvert) FromData(author *persistence.Author) *entity.Author {
	result := entity.Author{}
	_ = copier.Copy(&result, &author)
	return &result
}

func (a *AuthorConvert) ToData(author *entity.Author) *persistence.Author {
	result := persistence.Author{}
	_ = copier.Copy(&result, &author)
	return &result
}
