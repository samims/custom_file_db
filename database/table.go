package database

import (
	"custom_db/constants"
	"custom_db/wrapper"
	"fmt"
	"os"
	"strings"
)

type TableHandler struct {
	fileOperator wrapper.FileOperator
}

// NewTableHandler creates a new instance of TableHandler and returns a pointer to it.
// The TableHandler struct is used to handle operations related to a table, such as inserting values into it.
func NewTableHandler(fileOperator wrapper.FileOperator) *TableHandler {
	return &TableHandler{
		fileOperator: fileOperator,
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

	//_, err = fmt.Fprintf(file, "%s\n", strings.Join(values, ","))
	strToWrite := fmt.Sprintf("%s", strings.Join(values, ","))
	_, err = t.fileOperator.WriteString(file, strToWrite)
	if err != nil {
		return fmt.Errorf("error writing to table file: %w", err)
	}
	return nil
}
