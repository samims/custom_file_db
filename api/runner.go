package api

import (
	"context"
	"custom_db/database"
	"custom_db/sqlparser"
	"custom_db/wrapper"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
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
	config := NewConfig(viper.New())

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.GetRedisAddress(),
		Password: config.GetRedisPass(),   // no password set
		DB:       config.GetRedisDBName(), // use default DB
	})

	redisDB := wrapper.NewRedisOperator(redisClient)

	fileOperator := wrapper.NewFileOperator()
	metadataHandler := database.NewMetadataHandler(fileOperator)
	tableHandler := database.NewTableHandler(fileOperator, metadataHandler)

	sqlParser := sqlparser.NewSqlParser(metadataHandler, tableHandler, redisDB)

	router := Init(*sqlParser)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.GetAppPort()),
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
