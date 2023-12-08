package sqlparser

import (
	"custom_db/database"
	"fmt"
	"strings"
)

// SqlParser represents a SQL parser that handles metadata and table operations.
type SqlParser struct {
	MetadataHandler *database.MetadataHandler
	TableHandler    *database.TableHandler
}

// NewSqlParser creates a new instance of SqlParser initialized with the provided MetadataHandler and TableHandler.
// The MetadataHandler is used for handling metadata operations such as creating table metadata, while the TableHandler is used for handling table operations such as inserting data into
func NewSqlParser(metadataHandler *database.MetadataHandler, tableHandler *database.TableHandler) *SqlParser {
	return &SqlParser{
		MetadataHandler: metadataHandler,
		TableHandler:    tableHandler,
	}
}

func (s *SqlParser) ParseSQLQuery(query string) error {
	tokens := strings.Fields(query)
	if len(tokens) < 4 {
		return fmt.Errorf("invalid sql query")
	}
	if strings.ToUpper(tokens[0]) == "CREATE" && strings.ToUpper(tokens[1]) == "TABLE" {
		return s.handleCreateTable(tokens)
	}
	if strings.ToUpper(tokens[0]) == "INSERT" && strings.ToUpper(tokens[1]) == "INTO" {
		return s.handleInsertInto(tokens)
	}
	return fmt.Errorf("unsupported SQL operation")
}

func (s *SqlParser) handleCreateTable(tokens []string) error {
	colNames := strings.Split(tokens[3], ",")
	colTypes := strings.Split(tokens[4], ",")

	err := s.MetadataHandler.CreateTableMetadata(colNames, colTypes)

	if err != nil {
		return err
	}
	return nil

}

// handleInsertInto is responsible for implementing the logic to handle the "INSERT INTO" SQL operation.
// This method currently returns an error indicating that the insert operation is not yet implemented.
func (s *SqlParser) handleInsertInto(tokens []string) error {
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
	err := s.TableHandler.InsertIntoTable(extractedValues)
	if err != nil {
		return fmt.Errorf("error inserting values into table: %w", err)
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
