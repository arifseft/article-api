package domain

import (
	"context"
	"time"
)

// Article ...
type Article struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Body      string    `json:"body" validate:"required"`
	Author    string    `json:"author" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

// ArticleUsecase represent the article's usecases
type ArticleUsecase interface {
	Fetch(ctx context.Context, keyword string, author string) ([]Article, error)
	Store(context.Context, *Article) error
}

// ArticleRepository represent the article's repository contract
type ArticleRepository interface {
	Fetch(ctx context.Context, keyword string, author string) (res []Article, err error)
	Store(ctx context.Context, a *Article) error
}
