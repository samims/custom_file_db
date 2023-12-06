package database

import (
	"custom_db/constants"
	"fmt"
	"os"
	"strings"
)

// MetadataHandler is a struct that handles metadata files.
type MetadataHandler struct {
	MetadataFile string
}

// NewMetadataHandler creates a new instance of MetadataHandler with the
// provided metadataFile string.
func NewMetadataHandler(metadataFile string) *MetadataHandler {
	return &MetadataHandler{MetadataFile: metadataFile}
}

// CreateTableMetadata creates the metadata file for a new table with the given column names and column types.
// The method performs error checking to ensure that the metadata file does
func (m *MetadataHandler) CreateTableMetadata(colNames []string, colTypes []string) error {
	// error checking if table already exists the throw error,
	// this is only single table project
	if _, err := os.Stat(m.MetadataFile); err == nil {
		return fmt.Errorf("metadata file already exists")
	}

	if len(colNames) != len(colTypes) {
		return fmt.Errorf("number of column and type doesn't match")
	}

	file, err := os.OpenFile(m.MetadataFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening metadata file: %w", err)

	}
	defer file.Close()

	cols := make([]string, len(colNames))

	for i := range colNames {
		cols[i] = fmt.Sprintf("%s %s\n", colNames[i], colTypes[i])
	}
	_, err = fmt.Fprintf(file, "%s %s\n", constants.DefaultTableName, strings.Join(cols, ","))
	if err != nil {
		return fmt.Errorf("error writing to metadata file: %w", err)
	}
	return nil

}
