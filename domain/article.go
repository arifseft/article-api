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

type ArticleSearchPayload struct {
	Query  *string `json:"query"`
	Author *string `json:"author"`
}

// ArticleUsecase represent the article's usecases
type ArticleUsecase interface {
	GetArticles(ctx context.Context, payload ArticleSearchPayload) ([]Article, error)
	AddArticle(context.Context, *Article) error
}

// ArticleRepository represent the article's repository contract
type ArticleRepository interface {
	StoreArticle(ctx context.Context, a *Article) error
}

// ArticleSearch represent the article's search contract
type ArticleSearch interface {
	SearchArticle(ctx context.Context, payload ArticleSearchPayload) ([]Article, error)
	IndexArticle(ctx context.Context, a *Article) error
}

// ArticleEvent represent the article's event contract
type ArticleEvent interface {
	PublishArticleCreated(ctx context.Context, a Article) error
	SubscribeArticleCreated(ctx context.Context, f func(Article)) error
}

// ArticleCache represent the article's cache
type ArticleCache interface {
	GetCache(ctx context.Context, key string) (interface{}, error)
	SetCache(ctx context.Context, key string, value interface{}) error
	FlushAllCache(ctx context.Context) error
}
