package mysql_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	articleMysqlRepo "github.com/arifseft/article-api/article/repository/mysql"
	"github.com/arifseft/article-api/domain"
)

func TestFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockArticles := []domain.Article{
		domain.Article{
			ID: 1, Title: "article 1", Body: "body 1",
			Author: "John Wick", CreatedAt: time.Now(),
		},
		domain.Article{
			ID: 2, Title: "article 2", Body: "body 2",
			Author: "John Wick", CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body", "author", "created_at"}).
		AddRow(mockArticles[0].ID, mockArticles[0].Title, mockArticles[0].Body,
			mockArticles[0].Author, mockArticles[0].CreatedAt).
		AddRow(mockArticles[1].ID, mockArticles[1].Title, mockArticles[1].Body,
			mockArticles[1].Author, mockArticles[1].CreatedAt)

	keyword := "article"
	author := "John"
	query := "SELECT id, title, body, author, created_at FROM articles WHERE (title LIKE '%" + keyword + "%' OR body LIKE '%" + keyword + "%') AND author LIKE '%" + author + "%' ORDER BY created_at DESC"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := articleMysqlRepo.NewMysqlArticleRepository(db)
	list, err := a.Fetch(context.TODO(), keyword, author)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestStore(t *testing.T) {
	ar := &domain.Article{
		Title:     "Title",
		Body:      "Body",
		Author:    "M. Arif Sefrianto",
		CreatedAt: time.Now(),
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "INSERT articles SET title=\\? , body=\\? , author=\\?, created_at=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.Title, ar.Body, ar.Author, ar.CreatedAt).WillReturnResult(sqlmock.NewResult(12, 1))

	a := articleMysqlRepo.NewMysqlArticleRepository(db)

	err = a.Store(context.TODO(), ar)
	assert.NoError(t, err)
	assert.Equal(t, int64(12), ar.ID)
}
