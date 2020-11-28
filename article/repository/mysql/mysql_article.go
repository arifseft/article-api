package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/arifseft/article-api/domain"
	"github.com/sirupsen/logrus"
)

type mysqlArticleRepository struct {
	Conn *sql.DB
}

func NewMysqlArticleRepository(Conn *sql.DB) domain.ArticleRepository {
	return &mysqlArticleRepository{Conn}
}

func (m *mysqlArticleRepository) Fetch(ctx context.Context, query string, author string) (res []domain.Article, err error) {
	q := `SELECT id, title, body, author, created_at FROM articles `

	var where []string
	if query != "" {
		where = append(where, fmt.Sprintf("(title LIKE '%%%s%%' OR body LIKE '%%%s%%') ", query, query))
	}
	if author != "" {
		where = append(where, fmt.Sprintf("author LIKE '%%%s%%' ", author))
	}
	if len(where) > 0 {
		q += " WHERE " + strings.Join(where, " AND ")
	}
	q += ` ORDER BY created_at DESC`

	rows, err := m.Conn.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		t := domain.Article{}
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Body,
			&t.Author,
			&t.CreatedAt,
		)
		res = append(res, t)
	}

	return
}

func (m *mysqlArticleRepository) Store(ctx context.Context, a *domain.Article) (err error) {
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
