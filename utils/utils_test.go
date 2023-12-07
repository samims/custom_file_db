package utils

import (
	"testing"
)

func TestFormatMetadataString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		colNames []string
		colTypes []string
		expected string
	}{
		{
			name:     "empty slices",
			colNames: []string{},
			colTypes: []string{},
			expected: "",
		},
		{
			name:     "single element",
			colNames: []string{"id"},
			colTypes: []string{"int"},
			expected: "id int\n",
		},
		{
			name:     "multiple elements",
			colNames: []string{"id", "name"},
			colTypes: []string{"int", "string"},
			expected: "id int\nname string\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			result := FormatMetadataString(tc.colNames, tc.colTypes)
			if result != tc.expected {
				t.Errorf("Expected: %q, got: %q", tc.expected, result)
			}
		})
	}
}
