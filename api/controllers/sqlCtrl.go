package controllers

import (
	"custom_db/api/types"
	"custom_db/sqlparser"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SQLCtrl interface {
	Query(ctx *gin.Context)
}

type sqlCtrl struct {
	sqlParser sqlparser.SqlParser
}

func (s *sqlCtrl) Query(ctx *gin.Context) {
	var data types.RequestData
	err := ctx.BindJSON(&data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(data.Query)
	err = s.sqlParser.ParseSQLQuery(data.Query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Query executed successfully"})
}

func NewSQLCtrl(sqlParser sqlparser.SqlParser) SQLCtrl {
	return &sqlCtrl{
		sqlParser: sqlParser,
	}
}
