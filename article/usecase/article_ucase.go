package usecase

import (
	"context"
	"time"

	"github.com/arifseft/article-api/domain"
)

type articleUsecase struct {
	articleRepo    domain.ArticleRepository
	contextTimeout time.Duration
}

func NewArticleUsecase(a domain.ArticleRepository, timeout time.Duration) domain.ArticleUsecase {
	return &articleUsecase{
		articleRepo:    a,
		contextTimeout: timeout,
	}
}

func (a *articleUsecase) Fetch(c context.Context, keyword string, author string) (res []domain.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.articleRepo.Fetch(ctx, keyword, author)
	if err != nil {
		return nil, err
	}

	return
}

func (a *articleUsecase) Store(c context.Context, m *domain.Article) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	m.CreatedAt = time.Now()

	err = a.articleRepo.Store(ctx, m)
	return
}
