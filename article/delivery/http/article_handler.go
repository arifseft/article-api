package http

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/arifseft/article-api/domain"
)

type ArticleHandler struct {
	AUsecase domain.ArticleUsecase
}

func NewArticleHandler(e *echo.Echo, us domain.ArticleUsecase) {
	handler := &ArticleHandler{
		AUsecase: us,
	}
	e.GET("/articles", handler.FetchArticle)
	e.POST("/articles", handler.Store)
}

type Response struct {
	Data    interface{} `json:"data"`
	Message *string     `json:"message"`
}

func (a *ArticleHandler) FetchArticle(c echo.Context) error {
	query := c.QueryParam("query")
	author := c.QueryParam("author")
	ctx := c.Request().Context()

	listAr, err := a.AUsecase.Fetch(ctx, query, author)
	if err != nil {
		errStr := err.Error()
		return c.JSON(getStatusCode(err), Response{
			Data:    nil,
			Message: &errStr,
		})
	}

	return c.JSON(http.StatusOK, Response{
		Data:    listAr,
		Message: nil,
	})
}

func isRequestValid(m *domain.Article) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *ArticleHandler) Store(c echo.Context) (err error) {
	var article domain.Article
	err = c.Bind(&article)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&article); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = a.AUsecase.Store(ctx, &article)
	if err != nil {
		errStr := err.Error()
		return c.JSON(getStatusCode(err), Response{
			Data:    nil,
			Message: &errStr,
		})
	}

	return c.JSON(http.StatusCreated, article)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	return http.StatusInternalServerError
}
