package sqlparser

import (
	"custom_db/database/mocks"
	wrapperMocks "custom_db/wrapper/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqlParser(t *testing.T) {
	metadataHandlerMock := mocks.NewMetadataHandler(t)
	tableHandlerMock := mocks.NewTableHandler(t)
	redisMock := wrapperMocks.NewRedisOperator(t)

	parser := NewSqlParser(metadataHandlerMock, tableHandlerMock, redisMock)

	t.Run("parse empty query", func(t *testing.T) {
		_, err := parser.ParseSQLQuery("")
		assert.Equal(t, "invalid sql query", err.Error())
	})

	t.Run("parse unsupported operation", func(t *testing.T) {
		_, err := parser.ParseSQLQuery("INVALID * FROM table_name")
		assert.Equal(t, "unsupported SQL operation", err.Error())
	})

	t.Run("parse create table without enough tokens", func(t *testing.T) {
		_, err := parser.ParseSQLQuery("CREATE TABLE table_name")
		assert.Equal(t, "invalid sql query", err.Error())
	})

	t.Run("parse insert into without values", func(t *testing.T) {
		_, err := parser.ParseSQLQuery("INSERT INTO students values")
		assert.Equal(t, "no values found in insert query", err.Error())
	})

}

func TestExtractValues(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "extract without parentheses",
			input:    "value1, value2, value3",
			expected: []string{},
		},
		{
			name:     "extract single value",
			input:    "(value1)",
			expected: []string{"value1"},
		},
		{
			name:     "extract multiple values",
			input:    "(value1, value2, value3)",
			expected: []string{"value1, value2, value3"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp := extractValues(tt.input)
			assert.Equal(t, tt.expected, resp)
		})
	}
}
