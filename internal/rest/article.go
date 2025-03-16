package rest

import (
	"github.com/bxcodec/go-clean-arch/domain/domain/entity"
	"github.com/bxcodec/go-clean-arch/internal"
	"github.com/bxcodec/go-clean-arch/internal/application"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// ArticleHandler  represent the httphandler for article
type ArticleHandler struct {
	Service application.ArticleService
}

const defaultNum = 10

// NewArticleHandler will initialize the articles/ resources endpoint
func NewArticleHandler(e *echo.Echo, svc application.ArticleService) {
	handler := &ArticleHandler{
		Service: svc,
	}
	e.GET("/articles", handler.FetchArticle)
	e.POST("/articles", handler.Store)
	e.GET("/articles/:id", handler.GetByID)
	e.DELETE("/articles/:id", handler.Delete)
}

// FetchArticle will fetch the article based on given params
func (a *ArticleHandler) FetchArticle(c echo.Context) error {

	numS := c.QueryParam("num")
	num, err := strconv.Atoi(numS)
	if err != nil || num == 0 {
		num = defaultNum
	}

	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()

	listAr, nextCursor, err := a.Service.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, listAr)
}

// GetByID will get article by given id
func (a *ArticleHandler) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, internal.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	art, err := a.Service.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, art)
}

func isRequestValid(m *entity.Article) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the article by given request body
func (a *ArticleHandler) Store(c echo.Context) (err error) {
	var article entity.Article
	err = c.Bind(&article)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&article); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = a.Service.Store(ctx, &article)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, article)
}

// Delete will delete article by given param
func (a *ArticleHandler) Delete(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, internal.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	err = a.Service.Delete(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case internal.ErrInternalServerError:
		return http.StatusInternalServerError
	case internal.ErrNotFound:
		return http.StatusNotFound
	case internal.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
