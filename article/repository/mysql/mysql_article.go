package mysql

import (
	"context"
	"database/sql"

	"github.com/arifseft/article-api/domain"
)

type mysqlArticleRepository struct {
	Conn *sql.DB
}

func NewMysqlArticleRepository(Conn *sql.DB) domain.ArticleRepository {
	return &mysqlArticleRepository{Conn}
}

func (m *mysqlArticleRepository) StoreArticle(ctx context.Context, a *domain.Article) (err error) {
	q := `INSERT articles SET title=? , body=? , author=?, created_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, q)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, a.Title, a.Body, a.Author, a.CreatedAt)
	if err != nil {
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	a.ID = lastID
	return
}
