package main

import (
	"context"
	"log"
	"time"

	"github.com/arifseft/article-api/domain"
	"github.com/tinrab/retry"
)

func ConsumeArticleCreated(natsArticleEvent domain.ArticleEvent, elasticArticleRepository domain.ArticleRepository) (err error) {
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		err = natsArticleEvent.SubscribeArticleCreated(context.Background(), func(m domain.Article) {
			if err := elasticArticleRepository.Store(context.Background(), &m); err != nil {
				log.Println(err)
			}
		})
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	})
	return nil
}
