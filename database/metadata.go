package database

import (
	"bufio"
	"custom_db/utils"
	"custom_db/wrapper"
	"errors"
	"fmt"
	"os"
	"strings"
)

type MetadataHandler interface {
	CreateTableMetadata(tableName string, colNames []string, colTypes []string) error
	ReadColumnTypes(filename string) (map[string]string, error)
	ReadColNamesAndTypesInArray(fileName string) ([]string, []string, error)
}

// MetadataHandler is a struct that handles metadata files.
type metadataHandler struct {
	fileOperator wrapper.FileOperator
}

// NewMetadataHandler creates a new instance of MetadataHandler with the
// provided metadataFile string.
func NewMetadataHandler(fileOperator wrapper.FileOperator) MetadataHandler {
	return &metadataHandler{
		fileOperator: fileOperator,
	}
}

// CreateTableMetadata creates the metadata file for a new table with the given
// column names and column types. The method performs error checking
// to ensure that the metadata file does
func (m *metadataHandler) CreateTableMetadata(tableName string, colNames []string, colTypes []string) error {
	if !utils.IsDirEmpty() {
		return errors.New(`table already exists. please delete as this is a single table system`)
	}

	var metadataFileName = utils.GetMetadataFileName(tableName)
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

func (m *metadataHandler) ReadColumnTypes(filename string) (map[string]string, error) {
	file, err := m.fileOperator.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer m.fileOperator.CloseFile(file)

	metaTypes := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid metadata format: %s", line)
		}

		metaTypes[parts[0]] = parts[1]
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	return metaTypes, nil
}

func (m *metadataHandler) ReadColNamesAndTypesInArray(fileName string) ([]string, []string, error) {
	types, err := m.ReadColumnTypes(fileName)
	if err != nil {
		return nil, nil, err
	}

	var colNames []string
	var colTypes []string
	for name, typ := range types {
		colNames = append(colNames, name)
		colTypes = append(colTypes, typ)
	}
	return colNames, colTypes, nil
}
