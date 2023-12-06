package database

import (
	"fmt"
	"os"
	"testing"
)

func TestCreateTableMetadata(t *testing.T) {
	type fields struct {
		metaDataFileName string
	}

	type args struct {
		colNames []string
		colTypes []string
	}

	cleanUp := func(f *fields) {
		err := os.Remove(f.metaDataFileName)
		if err != nil {
			fmt.Println("error occurred during test metaData file remove")
		}
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantToPass bool
		beforeTest func(*fields)
		afterTest  func(*fields)
	}{
		{
			"ValidInput",
			fields{metaDataFileName: "test_file"},
			args{
				colNames: []string{"id", "name"},
				colTypes: []string{"int", "string"},
			},
			true,
			func(f *fields) {},
			cleanUp,
		},

		{
			"MetadataFileExists",
			fields{metaDataFileName: "test_file"},
			args{
				[]string{"id", "name"},
				[]string{"int", "string"},
			},
			false,
			func(f *fields) {
				_, err := os.Create(f.metaDataFileName)
				if err != nil {
					return
				}
			},
			cleanUp,
		},
		{
			"MismatchColumnTypes",
			fields{"testfile"},
			args{
				[]string{"id"},
				[]string{"int", "string"},
			},
			false,
			func(f *fields) {},
			cleanUp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeTest(&tt.fields)
			defer tt.afterTest(&tt.fields)

			mh := NewMetadataHandler(tt.fields.metaDataFileName)
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
