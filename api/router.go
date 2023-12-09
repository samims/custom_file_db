package api

import (
	"custom_db/api/controllers"
	"custom_db/sqlparser"
	"github.com/gin-gonic/gin"
)

func Init(sqlParser sqlparser.SqlParser) *gin.Engine {
	router := gin.Default()

	queryCtrl := controllers.NewSQLCtrl(sqlParser)

	router.POST("/query", queryCtrl.Query)
	return router
}
