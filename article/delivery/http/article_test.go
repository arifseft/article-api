package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	articleHttp "github.com/arifseft/article-api/article/delivery/http"
	"github.com/arifseft/article-api/domain"
	"github.com/arifseft/article-api/domain/mocks"
)

func TestFetch(t *testing.T) {
	var mockArticle domain.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)
	mockUCase := new(mocks.ArticleUsecase)
	mockListArticle := make([]domain.Article, 0)
	mockListArticle = append(mockListArticle, mockArticle)
	query := "Abc"
	author := "John"
	payload := domain.ArticleSearchPayload{
		Query:  &query,
		Author: &author,
	}
	mockUCase.On("GetArticles", context.TODO(), payload).Return(mockListArticle, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/article?query="+query+"&author="+author, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := articleHttp.ArticleHandler{
		AUsecase: mockUCase,
	}
	err = handler.FetchArticle(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestFetchError(t *testing.T) {
	mockUCase := new(mocks.ArticleUsecase)
	query := "1"
	author := "John"
	payload := domain.ArticleSearchPayload{
		Query:  &query,
		Author: &author,
	}
	mockUCase.On("GetArticles", context.TODO(), payload).Return(nil, domain.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/article?query=1&author="+author, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := articleHttp.ArticleHandler{
		AUsecase: mockUCase,
	}
	err = handler.FetchArticle(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestStore(t *testing.T) {
	mockArticle := domain.Article{
		Title:     "Title",
		Body:      "Body",
		Author:    "Author",
		CreatedAt: time.Now(),
	}

	tempMockArticle := mockArticle
	tempMockArticle.ID = 0
	mockUCase := new(mocks.ArticleUsecase)

	j, err := json.Marshal(tempMockArticle)
	assert.NoError(t, err)

	mockUCase.On("AddArticle", context.TODO(), mock.AnythingOfType("*domain.Article")).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/article", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/article")

	handler := articleHttp.ArticleHandler{
		AUsecase: mockUCase,
	}
	err = handler.Store(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUCase.AssertExpectations(t)
}
