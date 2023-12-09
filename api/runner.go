package api

import (
	"context"
	"custom_db/database"
	"custom_db/sqlparser"
	"custom_db/wrapper"
	"fmt"
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

	fileOperator := wrapper.NewFileOperator()
	metadataHandler := database.NewMetadataHandler(fileOperator)
	tableHandler := database.NewTableHandler(fileOperator)

	sqlParser := sqlparser.NewSqlParser(metadataHandler, tableHandler)

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
