package utils

import (
	"fmt"
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
