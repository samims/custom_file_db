package database

import (
	"custom_db/constants"
	"fmt"
	"os"
	"testing"
)

func TestCreateTableMetadata(t *testing.T) {

	type args struct {
		colNames []string
		colTypes []string
	}

	cleanUp := func() {
		err := os.Remove(constants.DefaultTableMetadataName + ".txt")
		if err != nil {
			fmt.Println("error occurred during test metaData file remove")
		}
	}

	tests := []struct {
		name       string
		args       args
		wantToPass bool
		beforeTest func()
		afterTest  func()
	}{
		{
			"ValidInput",
			args{
				colNames: []string{"id", "name"},
				colTypes: []string{"int", "string"},
			},
			true,
			func() {},
			cleanUp,
		},

		{
			"MetadataFileExists",

			args{
				[]string{"id", "name"},
				[]string{"int", "string"},
			},
			false,
			func() {
				_, err := os.Create(constants.DefaultTableMetadataName + ".txt")
				if err != nil {
					return
				}
			},
			cleanUp,
		},
		{
			"MismatchColumnTypes",

			args{
				[]string{"id"},
				[]string{"int", "string"},
			},
			false,
			func() {},
			func() {},
		},
		{
			"ErrorOpeningFile",
			args{
				[]string{"id"},
				[]string{"int", "string"},
			},
			false,
			func() {
			},
			func() {
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeTest()
			defer tt.afterTest()

			mh := NewMetadataHandler()
			err := mh.CreateTableMetadata(tt.args.colNames, tt.args.colTypes)

			if tt.wantToPass && err != nil {
				t.Errorf("CreateTableMetadata() Error: %v, want nil error", err)
			}

			if !tt.wantToPass && err == nil {
				t.Errorf("CreateTableMetadata() Error: nil, want error")
			}

		})
	}
}
