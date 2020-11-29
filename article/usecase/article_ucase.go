package usecase

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/arifseft/article-api/domain"
)

type articleUsecase struct {
	articleMysqlRepo   domain.ArticleRepository
	articleElasticRepo domain.ArticleRepository
	articleNatsEvent   domain.ArticleEvent
	articleRedisCache  domain.ArticleCache
	contextTimeout     time.Duration
}

func NewArticleUsecase(
	articleMysqlRepository domain.ArticleRepository,
	articleElasticRepository domain.ArticleRepository,
	articleNatsEvent domain.ArticleEvent,
	articleRedisCache domain.ArticleCache,
	timeout time.Duration,
) domain.ArticleUsecase {
	return &articleUsecase{
		articleMysqlRepo:   articleMysqlRepository,
		articleElasticRepo: articleElasticRepository,
		articleNatsEvent:   articleNatsEvent,
		articleRedisCache:  articleRedisCache,
		contextTimeout:     timeout,
	}
}

func (a *articleUsecase) Fetch(c context.Context, query string, author string) (res []domain.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	// Get cache
	var keys []string
	keys = append(keys, "article")
	if author != "" {
		keys = append(keys, author)
	}
	if query != "" {
		keys = append(keys, query)
	}
	cacheKey := strings.Join(keys, ":")
	cacheValue, err := a.articleRedisCache.GetCache(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	if cacheValue != nil {
		err = json.Unmarshal([]byte(cacheValue.(string)), &res)
		if err != nil {
			return nil, err
		}
	} else {
		res, err = a.articleElasticRepo.Fetch(ctx, query, author)
		if err != nil {
			return nil, err
		}

		body, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}

		// Set cache
		err = a.articleRedisCache.SetCache(ctx, cacheKey, string(body))
		if err != nil {
			return nil, err
		}
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
