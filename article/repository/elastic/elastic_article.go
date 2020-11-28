package elastic

import (
	"context"
	"log"
	"reflect"
	"strconv"

	"github.com/arifseft/article-api/domain"
	"github.com/olivere/elastic/v7"
)

type elasticArticleRepository struct {
	Client *elastic.Client
}

const (
	index = "article"
)

func NewElasticArticleRepository(Client *elastic.Client) domain.ArticleRepository {
	return &elasticArticleRepository{Client}
}

func (e *elasticArticleRepository) Fetch(ctx context.Context, query string, author string) (res []domain.Article, err error) {
	searchSource := elastic.NewSearchSource()

	if query != "" {
		searchSource.Query(elastic.NewMatchQuery("title", query))
		searchSource.Query(elastic.NewMatchQuery("body", query))
	}

	if author != "" {
		searchSource.Query(elastic.NewMatchQuery("author", author))
	}

	searchService := e.Client.Search().Index(index).SearchSource(searchSource)

	searchResult, err := searchService.Do(ctx)
	if err != nil {
		log.Println("[ProductsES][GetPIds]Error=", err)
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

func (e *elasticArticleRepository) Store(ctx context.Context, a *domain.Article) (err error) {
	exist, err := e.Client.IndexExists(index).Do(ctx)
	if err != nil {
		log.Fatalf("IndexExists() ERROR: %v", err)
		return

	} else if !exist {
		createdIndex := e.Client.CreateIndex(index)
		if createdIndex == nil {
			log.Fatalf("CreateIndex() ERROR: %v", err)
			return
		}
	}

	var id string = "article_" + strconv.Itoa(int(a.ID))

	_, err = e.Client.Index().
		Index(index).
		Type("_doc").
		BodyJson(a).
		Id(id).
		Do(ctx)

	if err != nil {
		log.Fatalf("client.Index() ERROR: %v", err)
	}

	return
}
