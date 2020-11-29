package redis

import (
	"context"
	"time"

	"github.com/arifseft/article-api/domain"
	"github.com/go-redis/redis"
)

const (
	expires = 5 * time.Minute
)

type redisArticleCache struct {
	Client *redis.Client
}

func NewRedisCache(Client *redis.Client) domain.ArticleCache {
	return &redisArticleCache{Client: Client}
}

func (r *redisArticleCache) GetCache(ctx context.Context, key string) (value interface{}, err error) {
	value, err = r.Client.Get(key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return
}

func (r *redisArticleCache) SetCache(ctx context.Context, key string, value interface{}) (err error) {
	err = r.Client.Set(key, value, expires).Err()
	if err != nil {
		return err
	}
	return
}
