package database

import (
	"custom_db/constants"
	"custom_db/wrapper"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TableHandler struct {
	fileOperator    wrapper.FileOperator
	metaDataHandler MetadataHandler
}

// NewTableHandler creates a new instance of TableHandler and returns a pointer to it.
// The TableHandler struct is used to handle operations related to a table, such as inserting values into it.
func NewTableHandler(fileOperator wrapper.FileOperator, metadataHandler MetadataHandler) *TableHandler {
	return &TableHandler{
		fileOperator:    fileOperator,
		metaDataHandler: metadataHandler,
	}
}

// InsertIntoTable opens the table file in append mode, writes the values to it separated by commas, and closes the file.
// Parameters:
// - values: a slice of strings representing the values to be inserted into the table
// Returns:
// - error: if there is an error opening or writing to the table file
func (t *TableHandler) InsertIntoTable(values []string) error {
	//file, err := os.OpenFile(constants.DefaultTableName+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	file, err := t.fileOperator.OpenFile(constants.DefaultTableName+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening table file %w", err)
	}
	//defer file.Close()
	defer t.fileOperator.CloseFile(file)
	err = t.ValidateDataType(values)
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

func (t *TableHandler) ValidateDataType(values []string) error {
	metadataFile := constants.DefaultTableMetadataName + ".txt"
	types, err := t.metaDataHandler.ReadColumnTypes(metadataFile)
	if err != nil {
		return err
	}

	return t.validateColumnTypes(types, values)
}

func (t *TableHandler) validateColumnTypes(types map[string]string, values []string) error {
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

func validateColumnType(colName, colType, value string) error {
	// Validation logic for each data type (int, string, float, etc.)
	// Implement the respective checks for the data type
	// Example:
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
