package api

import (
	"context"
	"custom_db/database"
	"custom_db/sqlparser"
	"custom_db/wrapper"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
	"sync"
	"time"
)

type APIRunner interface {
	Go(ctx context.Context, wg *sync.WaitGroup)
}

type apiRunner struct {
}

func (runner *apiRunner) Go(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	redisDB := wrapper.NewRedisOperator(redisClient)

	fileOperator := wrapper.NewFileOperator()
	metadataHandler := database.NewMetadataHandler(fileOperator)
	tableHandler := database.NewTableHandler(fileOperator, metadataHandler)

	sqlParser := sqlparser.NewSqlParser(metadataHandler, tableHandler, redisDB)

	router := Init(*sqlParser)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8080),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
		return
	}
}

func NewAPIRunner() APIRunner {
	return &apiRunner{}
}
