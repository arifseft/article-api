package usecase

import (
	"context"
	"time"

	"github.com/arifseft/article-api/domain"
)

type articleUsecase struct {
	articleMysqlRepo   domain.ArticleRepository
	articleElasticRepo domain.ArticleRepository
	articleNatsEvent   domain.ArticleEvent
	contextTimeout     time.Duration
}

func NewArticleUsecase(
	articleMysqlRepository domain.ArticleRepository,
	articleElasticRepository domain.ArticleRepository,
	articleNatsEvent domain.ArticleEvent,
	timeout time.Duration,
) domain.ArticleUsecase {
	return &articleUsecase{
		articleMysqlRepo:   articleMysqlRepository,
		articleElasticRepo: articleElasticRepository,
		articleNatsEvent:   articleNatsEvent,
		contextTimeout:     timeout,
	}
}

func (a *articleUsecase) Fetch(c context.Context, keyword string, author string) (res []domain.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.articleElasticRepo.Fetch(ctx, keyword, author)
	if err != nil {
		return nil, err
	}

	return
}

func (a *articleUsecase) Store(c context.Context, m *domain.Article) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	m.CreatedAt = time.Now()

	err = a.articleMysqlRepo.Store(ctx, m)

	// Publish event
	err = a.articleNatsEvent.PublishArticleCreated(ctx, *m)
	return
}
