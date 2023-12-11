package sqlparser

import (
	"context"
	"custom_db/constants"
	"custom_db/database"
	"custom_db/utils"
	"custom_db/wrapper"
	"errors"
	"fmt"
	"strings"
	"time"
)

// SqlParser represents a SQL parser that handles metadata and table operations.
type SqlParser struct {
	MetadataHandler database.MetadataHandler
	TableHandler    database.TableHandler
	RedisDB         wrapper.RedisOperator
}

// NewSqlParser creates a new instance of SqlParser initialized with the provided MetadataHandler and TableHandler.
// The MetadataHandler is used for handling metadata operations such as creating table metadata, while the TableHandler is used for handling table operations such as inserting data into
func NewSqlParser(metadataHandler database.MetadataHandler, tableHandler database.TableHandler, redisDB wrapper.RedisOperator) *SqlParser {
	return &SqlParser{
		MetadataHandler: metadataHandler,
		TableHandler:    tableHandler,
		RedisDB:         redisDB,
	}
}

func (s *SqlParser) ParseSQLQuery(query string) ([]map[string]any, error) {
	var result []map[string]any
	var err error
	successfulExecution := true

	tokens := strings.Fields(query)
	if len(tokens) < 4 {
		if !(len(tokens) == 3 && strings.ToUpper(tokens[0]) == "DROP") {
			err = fmt.Errorf("invalid sql query")
			pushLogToRedis(s.RedisDB, query, false)
			return result, err
		}
	}
	if strings.ToUpper(tokens[0]) == "CREATE" && strings.ToUpper(tokens[1]) == "TABLE" {
		tableName := tokens[2]
		err = s.handleCreateTable(tableName, query)
		if err != nil {
			successfulExecution = false
		}
		pushLogToRedis(s.RedisDB, query, successfulExecution)
		return result, err
	}
	// for Drop table
	if strings.ToUpper(tokens[0]) == "DROP" && strings.ToUpper(tokens[1]) == "TABLE" {
		tableName := tokens[2]
		err = s.handleDropTable(tableName)
		successfulExecution = err == nil
		pushLogToRedis(s.RedisDB, query, successfulExecution)
		return result, err
	}

	if strings.ToUpper(tokens[0]) == "INSERT" && strings.ToUpper(tokens[1]) == "INTO" {
		tableName := tokens[2]
		err = s.handleInsertInto(tableName, tokens)
		successfulExecution = err == nil
		pushLogToRedis(s.RedisDB, query, successfulExecution)
		return result, err
	}

	if strings.ToUpper(tokens[0]) == "SELECT" &&
		(strings.ToUpper(tokens[1]) == "*" || strings.ToUpper(tokens[1]) == "") &&
		strings.ToUpper(tokens[2]) == "FROM" {
		tableName := tokens[3]
		result, err := s.handleSelectFrom(tableName, query)
		successfulExecution = err == nil
		pushLogToRedis(s.RedisDB, query, successfulExecution)
		return result, err
	}

	pushLogToRedis(s.RedisDB, query, false)
	return result, fmt.Errorf("unsupported SQL operation")
}

// func (s *SqlParser) handleCreateTable(tokens []string) error {
func (s *SqlParser) handleCreateTable(tableName, query string) error {

	colNames, colTypes, err := extractColumns(query)
	if err != nil {
		return err
	}
	err = s.MetadataHandler.CreateTableMetadata(tableName, colNames, colTypes)
	if err != nil {
		return err
	}
	err = s.TableHandler.CreateEmptyTable(tableName)
	return err

}

// handleInsertInto is responsible for implementing the logic to handle the "INSERT INTO" SQL operation.
// This method currently returns an error indicating that the insert operation is not yet implemented.
func (s *SqlParser) handleInsertInto(tableName string, tokens []string) error {
	if len(tokens) < 4 {
		return fmt.Errorf("invalid `insert` query")
	}
	valuesPortion := strings.Join(tokens[4:], " ")

	// Extract values within parentheses
	extractedValues := extractValues(valuesPortion)
	if len(extractedValues) == 0 {
		return fmt.Errorf("no values found in insert query")
	}

	// insert into the table
	err := s.TableHandler.InsertIntoTable(tableName, extractedValues)

	if err != nil {
		return fmt.Errorf("error inserting values into table: %w", err)
	}
	return nil
}

func (s *SqlParser) handleSelectFrom(tableName, query string) ([]map[string]any, error) {
	metadataDataFile := utils.GetMetadataFileName(tableName)
	result := make([]map[string]any, 0)

	colNames, colTypes, err := s.MetadataHandler.ReadColNamesAndTypesInArray(metadataDataFile)
	if err != nil {
		return result, err
	}

	result, err = s.TableHandler.SelectFrom(tableName, query, colNames, colTypes)
	if err != nil {
		return nil, err
	}

	return result, err

}

func (s *SqlParser) handleDropTable(tableName string) error {
	err := s.TableHandler.DropTable(tableName)
	if err != nil {
		return fmt.Errorf("error dropping table: %w", err)
	}
	return nil

}

// extractValues extracts the values within parentheses from the provided string and returns them as a slice of strings.
// The function iterates over each character in the string and checks if it is '(' or ')'.
// If it is '(', it sets the "inParentheses" flag to true, indicating that the subsequent characters should be considered
// as part of the values.
// If it is ')', it set the "inParentheses" flag to false and appends the current value (trimmed) to the "values" slice.
func extractValues(str string) []string {
	var values = []string{}
	inParentheses := false

	var currentValue strings.Builder
	for _, char := range str {
		switch char {
		case '(':
			inParentheses = true
		case ')':
			inParentheses = false
			values = append(values, strings.TrimSpace(currentValue.String()))
		default:
			if inParentheses {
				currentValue.WriteRune(char)
			}

		}

	}
	return values
}

func extractColumns(query string) ([]string, []string, error) {
	// Split the string on white space
	splitQuery := strings.Split(query, " ")
	// Ensure the structure of the string is as expected
	if len(splitQuery) < 3 {
		return nil, nil, errors.New(constants.ErrQueryUnSupportedFormat)
	}

	tableDefinition := splitQuery[3:]

	definitionStr := strings.Join(tableDefinition, " ")
	definitionStr = strings.TrimPrefix(definitionStr, "(")
	definitionStr = strings.TrimSuffix(definitionStr, ")")

	var colNames []string
	var colTypes []string

	parts := strings.Split(definitionStr, ",")
	for _, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		subParts := strings.Fields(trimmedPart)
		if len(subParts) != 2 {
			return nil, nil, errors.New(constants.ErrInvalidColumnDefFormat)
		}
		cName := subParts[0]
		cType := subParts[1]

		// validate cType by checking it's only supported types
		if !isValidColumnType(cType) {
			return nil, nil, errors.New(constants.ErrUnSupportedColumnType)
		}

		colNames = append(colNames, cName)
		colTypes = append(colTypes, cType)
	}

	return colNames, colTypes, nil
}

// isValidColumnType checks if the given column type is valid.
// It returns true if the column type is either "int" or "string", false otherwise.
func isValidColumnType(cType string) bool {
	return (strings.ToLower(cType) == constants.IntegerType) || (strings.ToLower(cType) == constants.StringType)
}

func pushLogToRedis(rdb wrapper.RedisOperator, log string, status bool) {
	ctx := context.Background()

	myLog := wrapper.LogObject{
		Sql:       log,
		Status:    status,
		Timestamp: time.Now(),
	}

	go func(rdb wrapper.RedisOperator, myLog wrapper.LogObject) {
		err := rdb.Push(ctx, constants.RedisLogName, myLog)
		if err != nil {
			fmt.Println(err)
		}
	}(rdb, myLog)
}
