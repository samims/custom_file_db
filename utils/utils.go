package utils

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// FormatMetadataString formats the given column names and types into a
// metadata string with the format "name type\n". It returns the formatted
// string.
func FormatMetadataString(colNames, colTypes []string) string {
	cols := make([]string, len(colNames))
	for i := range colNames {
		cols[i] = fmt.Sprintf("%s %s\n", colNames[i], colTypes[i])
	}
	formattedStr := fmt.Sprintf("%s", strings.Join(cols, ""))
	return formattedStr
}

func GetMetadataFileName(tableName string) string {
	return fmt.Sprintf("data/%s.metadata", tableName)
}

func GetTableFileName(tableName string) string {
	return fmt.Sprintf("data/%s.txt", tableName)
}

func IsDirEmpty() bool {
	// check if /data directory empty
	// create directory if not exists
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	fi, err := os.Open("data")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer fi.Close()
	_, err = fi.Readdirnames(1)
	if err == io.EOF {
		return true
	}
	return false
}
