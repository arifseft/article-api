package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/nats-io/nats.go"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"

	_articleRedisCache "github.com/arifseft/article-api/article/cache/redis"
	_articleHttpDelivery "github.com/arifseft/article-api/article/delivery/http"
	_articleHttpDeliveryMiddleware "github.com/arifseft/article-api/article/delivery/http/middleware"
	_articleNatsEvent "github.com/arifseft/article-api/article/event/nats"
	_articleMysqlRepo "github.com/arifseft/article-api/article/repository/mysql"
	_articleElasticSearch "github.com/arifseft/article-api/article/search/elastic"
	_articleUcase "github.com/arifseft/article-api/article/usecase"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	var err error
	// Connect to MySQL
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to ElasticSearch
	esHost := viper.GetString(`elasticsearch.host`)
	elasticClient, err := elastic.NewClient(
		elastic.SetSniff(true),
		elastic.SetURL(esHost),
		elastic.SetHealthcheckInterval(5*time.Second), // quit trying after 5 seconds
	)

	if err != nil {
		log.Fatal(err)
	}

	// Connect to Redis
	rdClient := redis.NewClient(&redis.Options{
		Addr:     viper.GetString(`redis.host`),
		DB:       viper.GetInt(`redis.db`),
		Password: viper.GetString(`redis.pass`),
	})

	// Connect to NATS
	natsHost := viper.GetString(`nats.host`)
	natsConn, err := nats.Connect(natsHost)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
		err = rdClient.Close()
		if err != nil {
			log.Fatal(err)
		}
		natsConn.Close()
	}()

	e := echo.New()
	middL := _articleHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	mysqlArticleRepository := _articleMysqlRepo.NewMysqlArticleRepository(dbConn)
	elasticArticleSearch := _articleElasticSearch.NewElasticArticleSearch(elasticClient)
	natsArticleEvent := _articleNatsEvent.NewNatsArticleEvent(natsConn)
	redisArticleCache := _articleRedisCache.NewRedisCache(rdClient)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	au := _articleUcase.NewArticleUsecase(mysqlArticleRepository, elasticArticleSearch, natsArticleEvent, redisArticleCache, timeoutContext)
	_articleHttpDelivery.NewArticleHandler(e, au)

	ConsumeArticleCreated(natsArticleEvent, elasticArticleSearch)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
