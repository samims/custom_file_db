package database

import (
	"custom_db/constants"
	"custom_db/utils"
	"custom_db/wrapper"
	"fmt"
	"os"
)

// MetadataHandler is a struct that handles metadata files.
type MetadataHandler struct {
	fileOperator wrapper.FileOperator
}

// NewMetadataHandler creates a new instance of MetadataHandler with the
// provided metadataFile string.
func NewMetadataHandler(fileOperator wrapper.FileOperator) *MetadataHandler {
	return &MetadataHandler{
		fileOperator: fileOperator,
	}
}

// CreateTableMetadata creates the metadata file for a new table with the given column names and column types.
// The method performs error checking to ensure that the metadata file does
func (m *MetadataHandler) CreateTableMetadata(colNames []string, colTypes []string) error {
	// error checking if table already exists the throw error,
	// this is only single table project
	var metadataFileName = constants.DefaultTableMetadataName + ".txt"
	if _, err := m.fileOperator.Stat(metadataFileName); err == nil {
		return fmt.Errorf("metadata file already exists")
	}

	if len(colNames) != len(colTypes) {
		return fmt.Errorf("number of column and type doesn't match")
	}

	file, err := m.fileOperator.OpenFile(metadataFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error creating metadata file: %w", err)

	}

	defer m.fileOperator.CloseFile(file)

	formattedStr := utils.FormatMetadataString(colNames, colTypes)
	_, err = m.fileOperator.WriteString(file, formattedStr)
	if err != nil {
		return fmt.Errorf("error writing to metadata file: %w", err)
	}
	return nil
}
