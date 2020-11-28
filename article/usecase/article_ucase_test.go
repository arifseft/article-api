package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	ucase "github.com/arifseft/article-api/article/usecase"
	"github.com/arifseft/article-api/domain"
	"github.com/arifseft/article-api/domain/mocks"
)

func TestFetch(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := domain.Article{
		Title:  "Hello",
		Body:   "Abc",
		Author: "John",
	}

	mockListArticle := make([]domain.Article, 0)
	mockListArticle = append(mockListArticle, mockArticle)

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("Fetch", context.TODO(), mock.AnythingOfType("string"),
			mock.AnythingOfType("string")).Return(mockListArticle, nil).Once()

		u := ucase.NewArticleUsecase(mockArticleRepo, time.Second*2)
		author := "John"
		query := "Abc"
		list, err := u.Fetch(context.TODO(), query, author)
		assert.NotEmpty(t, list)
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListArticle))

		mockArticleRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepo.On("Fetch", context.TODO(), mock.AnythingOfType("string"),
			mock.AnythingOfType("string")).Return(nil, "", errors.New("Unexpexted Error")).Once()

		u := ucase.NewArticleUsecase(mockArticleRepo, time.Second*2)
		author := "John"
		query := "Abc"
		list, err := u.Fetch(context.TODO(), query, author)

		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockArticleRepo.AssertExpectations(t)
	})

}

func TestStore(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := domain.Article{
		Title:  "Hello",
		Body:   "Body",
		Author: "Author",
	}

	t.Run("success", func(t *testing.T) {
		tempMockArticle := mockArticle
		tempMockArticle.ID = 0
		mockArticleRepo.On("Store", context.TODO(), mock.AnythingOfType("*domain.Article")).Return(nil).Once()

		u := ucase.NewArticleUsecase(mockArticleRepo, time.Second*2)

		err := u.Store(context.TODO(), &tempMockArticle)

		assert.NoError(t, err)
		assert.Equal(t, mockArticle.Title, tempMockArticle.Title)
		mockArticleRepo.AssertExpectations(t)
	})
}
