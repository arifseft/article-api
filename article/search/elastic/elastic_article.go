package search

import (
	"context"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/arifseft/article-api/domain"
	"github.com/olivere/elastic/v7"
)

const (
	index = "article"
)

type elasticArticleSearch struct {
	Client *elastic.Client
}

func NewElasticArticleSearch(Client *elastic.Client) domain.ArticleSearch {
	return &elasticArticleSearch{Client}
}

func (e *elasticArticleSearch) SearchArticle(ctx context.Context, payload domain.ArticleSearchPayload) (res []domain.Article, err error) {
	if err := e.indexCheck(ctx, index); err != nil {
		return res, err
	}

	var shouldQuery []elastic.Query

	if payload.Query != nil {
		var queryQ *elastic.BoolQuery
		wildcard := "*" + strings.ToLower(*payload.Query) + "*"
		queryQ = elastic.NewBoolQuery()
		queryQ.Should(elastic.NewWildcardQuery("title", wildcard))
		queryQ.Should(elastic.NewWildcardQuery("body", wildcard))

		shouldQuery = append(shouldQuery, queryQ)
	}

	if payload.Author != nil {
		var authorQ *elastic.BoolQuery
		authorQ = elastic.NewBoolQuery()
		authorQ.Must(elastic.NewMatchQuery("author", strings.ToLower(*payload.Author)))

		shouldQuery = append(shouldQuery, authorQ)
	}

	newBoolQuery := elastic.NewBoolQuery().Must(shouldQuery...)

	searchResult, err := e.Client.Search().
		Index(index).
		Query(newBoolQuery).
		SortBy(elastic.NewFieldSort("created_at").Desc()).
		Do(ctx)

	if err != nil {
		log.Printf("SearchSource() ERROR: %v", err)
		return
	}

	var article domain.Article
	for _, item := range searchResult.Each(reflect.TypeOf(article)) {
		if t, ok := item.(domain.Article); ok {
			res = append(res, t)
		}
	}

	return
}

func (e *elasticArticleSearch) IndexArticle(ctx context.Context, a *domain.Article) (err error) {
	e.indexCheck(ctx, index)

	var id string = "article_" + strconv.Itoa(int(a.ID))

	_, err = e.Client.Index().
		Index(index).
		Type("_doc").
		BodyJson(a).
		Id(id).
		Do(ctx)

	if err != nil {
		log.Printf("Store() ERROR: %v", err)
	}

	return
}

func (e *elasticArticleSearch) indexCheck(ctx context.Context, index string) error {
	exist, err := e.Client.IndexExists(index).Do(ctx)
	if err != nil {
		log.Printf("IndexExists() ERROR: %v", err)
		return err
	}

	if !exist {
		createdIndex, err := e.Client.CreateIndex(index).Do(ctx)
		if err != nil {
			log.Printf("CreateIndex() ERROR: %v", err)
			return err
		}
		if createdIndex == nil {
			log.Printf("CreateIndex() ERROR: %v", err)
			return err
		}
	}
	return nil
}
