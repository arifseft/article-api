package elastic

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/arifseft/article-api/domain"
	"github.com/sirupsen/logrus"
)

type elasticArticleRepository struct {
	Conn *sql.DB
}

func NewElasticArticleRepository(Conn *sql.DB) domain.ArticleRepository {
	return &elasticArticleRepository{Conn}
}

func (e *elasticArticleRepository) Fetch(ctx context.Context, keyword string, author string) (res []domain.Article, err error) {
	query := `SELECT id, title, body, author, created_at FROM articles `

	var whereQuery []string
	if keyword != "" {
		whereQuery = append(whereQuery, fmt.Sprintf("(title LIKE '%%%s%%' OR body LIKE '%%%s%%') ", keyword, keyword))
	}
	if author != "" {
		whereQuery = append(whereQuery, fmt.Sprintf("author LIKE '%%%s%%' ", author))
	}
	if len(whereQuery) > 0 {
		query += " WHERE " + strings.Join(whereQuery, " AND ")
	}
	query += ` ORDER BY created_at DESC`
	// fmt.Printf("\nValue of  ========================================== : %v\n", )

	// res, err = e.fetch(ctx, query, keyword, keyword, author)

	rows, err := e.Conn.QueryContext(ctx, query)
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

func (e *elasticArticleRepository) Store(ctx context.Context, a *domain.Article) (err error) {
	query := `INSERT articles SET title=? , body=? , author=?, created_at=?`
	stmt, err := e.Conn.PrepareContext(ctx, query)
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
