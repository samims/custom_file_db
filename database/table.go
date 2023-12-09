package database

import (
	"bufio"
	"custom_db/constants"
	"custom_db/utils"
	"custom_db/wrapper"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TableHandler interface {
	CreateEmptyTable(tableName string) error
	InsertIntoTable(tableName string, values []string) error
	SelectFrom(tableName, query string, colNames, colTypes []string) ([]map[string]any, error)
	ValidateDataType(tableName string, values []string) error
	DropTable(tableName string) error
}

type tableHandler struct {
	fileOperator    wrapper.FileOperator
	metaDataHandler MetadataHandler
}

// NewTableHandler creates a new instance of TableHandler and returns a pointer to it.
// The TableHandler struct is used to handle operations related to a table, such as inserting values into it.
func NewTableHandler(fileOperator wrapper.FileOperator, metadataHandler MetadataHandler) TableHandler {
	return &tableHandler{
		fileOperator:    fileOperator,
		metaDataHandler: metadataHandler,
	}
}

func (t *tableHandler) CreateEmptyTable(tableName string) error {
	tableFileName := utils.GetTableFileName(tableName)
	// create empty file
	_, err := t.fileOperator.OpenFile(tableFileName, os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error creating table file: %w", err)
	}
	return nil
}

// InsertIntoTable opens the table file in append mode, writes the values to it separated by commas, and closes the file.
// Parameters:
// - values: a slice of strings representing the values to be inserted into the table
// Returns:
// - error: if there is an error opening or writing to the table file
func (t *tableHandler) InsertIntoTable(tableName string, values []string) error {
	//file, err := os.OpenFile(constants.DefaultTableName+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	tableFileName := utils.GetTableFileName(tableName)
	file, err := t.fileOperator.OpenFile(tableFileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening table file %w", err)
	}
	//defer file.Close()
	defer t.fileOperator.CloseFile(file)
	err = t.ValidateDataType(tableName, values)
	if err != nil {
		return err
	}
	strToWrite := fmt.Sprintf("%s", strings.Join(values, ","))
	_, err = t.fileOperator.WriteString(file, strToWrite)
	if err != nil {
		return fmt.Errorf("error writing to table file: %w", err)
	}
	return nil
}

func (t *tableHandler) ValidateDataType(tableName string, values []string) error {
	metadataFile := utils.GetMetadataFileName(tableName)
	types, err := t.metaDataHandler.ReadColumnTypes(metadataFile)
	if err != nil {
		return err
	}

	return t.validateColumnTypes(types, values)
}

func (t *tableHandler) validateColumnTypes(types map[string]string, values []string) error {
	index := 0
	valueStr := values[0]
	splitValueStr := strings.Split(valueStr, ",")
	fmt.Println(values[0])
	for colName, colType := range types {
		err := validateColumnType(colName, colType, splitValueStr[index])
		if err != nil {
			return err
		}
		index++
	}

	return nil
}

// SelectFrom opens the table file in read-only mode and reads the data from
// it based on the provided query. It returns a slice of maps, where each map
// represents a row from the table
func (t *tableHandler) SelectFrom(tableName, query string, colNames, colTypes []string) ([]map[string]any, error) {
	tableFileName := utils.GetTableFileName(tableName)
	file, err := t.fileOperator.OpenFile(tableFileName, os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("error opening table file: %w", err)
	}
	defer t.fileOperator.CloseFile(file)

	var result []map[string]any
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		row := make(map[string]any)
		for i, value := range parts {
			columnName := colNames[i]
			colType := colTypes[i]
			if strings.ToLower(colType) == constants.IntegerType {
				parsedVal, err := strconv.Atoi(value)
				if err != nil {
					return result, fmt.Errorf("invalid data type for column '%s': expected int", columnName)
				}
				row[columnName] = parsedVal
				continue

			} else if strings.ToLower(colType) == constants.StringType {
			} else {
				return result, fmt.Errorf("unsupported data type for column '%s': %s", columnName, colType)
			}
			row[columnName] = value
		}
		result = append(result, row)
	}

	// [{'id: 1, 'name' 'val1']
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning table file: %w", err)
	}

	return result, nil

}

func validateColumnType(colName, colType, value string) error {
	// Validation logic for each data type (int, string, float, etc.)
	// Implement the respective checks for the data type
	if strings.ToLower(colType) == "int" {
		if _, err := strconv.Atoi(value); err != nil {
			return fmt.Errorf("invalid data type for column '%s': expected int", colName)
		}
	} else if strings.ToLower(colType) == "string" {
		// No parsing needed for string type
		// parse to string
		fmt.Println("string type parsed")
	} else {
		return fmt.Errorf("unsupported data type for column '%s': %s", colName, colType)
	}
	return nil
}

// DropTable checks if the metadata and table files exist for the specified table name.
// If both files exist, it deletes them.
// If either file does not exist, it returns an error indicating that the table does not exist.
//
// Parameters:
// - tableName: a string representing the name of the table to drop
//
// Returns:
// - error: if there is an error checking or deleting the metadata or table file.
func (t *tableHandler) DropTable(tableName string) error {
	metadataFileName := utils.GetMetadataFileName(tableName)
	_, err := t.fileOperator.Stat(metadataFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("table '%s' does not exist", tableName)
		}
		return fmt.Errorf("error checking metadata file: %w", err)
	}
	tableFileName := utils.GetTableFileName(tableName)
	_, err = t.fileOperator.Stat(tableFileName)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("table '%s' does not exists", tableName)
		}
		return fmt.Errorf("error checking table file: %w", err)
	}

	err = t.fileOperator.DeleteFile(metadataFileName)
	if err != nil {
		return err
	}
	err = t.fileOperator.DeleteFile(tableFileName)
	if err != nil {
		return err
	}

	return nil

}
