package database

import (
	"custom_db/utils"
	"custom_db/wrapper"
	wrapperMocks "custom_db/wrapper/mocks"
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

func TestCreateTableMetadata(t *testing.T) {

	type fields struct {
		fileOperator *wrapperMocks.FileOperator
	}

	type args struct {
		tableName string
		colNames  []string
		colTypes  []string
	}

	cleanUp := func(f *fields, a *args) {
		metadataFileName := utils.GetMetadataFileName(a.tableName)
		err := os.Remove(metadataFileName)
		if err != nil {
			fmt.Println("error occurred during test metaData file remove")
		}
		tableFileName := utils.GetTableFileName(a.tableName)
		err = os.Remove(tableFileName)
		if err != nil {
			fmt.Println("error occurred during test table file remove")
		}
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantToPass bool
		beforeTest func(*fields, *args)
		afterTest  func(*fields, *args)
	}{
		{
			"ValidInput",
			fields{},
			args{
				colNames: []string{"id", "name"},
				colTypes: []string{"int", "string"},
			},
			true,
			func(f *fields, a *args) {
				//fileInfo :=
				f.fileOperator.On(
					"Stat",
					mock.Anything,
				).Return(&wrapper.MockFileInfo{}, errors.New("test"))
				var file os.File
				f.fileOperator.On(
					"OpenFile",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(&file, nil)

				f.fileOperator.On(
					"WriteString",
					mock.Anything,
					mock.Anything,
				).Return(int(0), nil)

				f.fileOperator.On(
					"CloseFile",
					mock.Anything,
				).Return(nil)
			},
			cleanUp,
		},

		{
			"MetadataFileExists",
			fields{},
			args{
				"abc",
				[]string{"id", "name"},
				[]string{"int", "string"},
			},
			false,
			func(f *fields, a *args) {
				//fileInfo :=
				f.fileOperator.On(
					"Stat",
					mock.Anything,
				).Return(&wrapper.MockFileInfo{}, nil)
			},
			cleanUp,
		},
		{
			"MismatchColumnTypeAndNameLength",
			fields{},
			args{
				"abc",
				[]string{"id"},
				[]string{"int", "string"},
			},
			false,
			func(f *fields, a *args) {
				//fileInfo :=
				f.fileOperator.On(
					"Stat",
					mock.Anything,
				).Return(&wrapper.MockFileInfo{}, errors.New("test"))
			},
			func(f *fields, a *args) {},
		},
		{
			"ErrorOpeningFile",
			fields{},
			args{
				"abc",
				[]string{"id"},
				[]string{"int"},
			},
			false,
			func(f *fields, a *args) {
				//fileInfo :=
				f.fileOperator.On(
					"Stat",
					mock.Anything,
				).Return(&wrapper.MockFileInfo{}, errors.New("test"))
				var file os.File
				f.fileOperator.On(
					"OpenFile",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(&file, errors.New("error creating metadata file"))
			},
			func(f *fields, a *args) {
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.fields = fields{
				fileOperator: wrapperMocks.NewFileOperator(t),
			}

			if tt.beforeTest != nil {
				tt.beforeTest(&tt.fields, &tt.args)
			}
			defer tt.afterTest(&tt.fields, &tt.args)

			mh := NewMetadataHandler(tt.fields.fileOperator)
			err := mh.CreateTableMetadata(tt.args.tableName, tt.args.colNames, tt.args.colTypes)

			if tt.wantToPass && err != nil {
				t.Errorf("CreateTableMetadata() Error: %v, want nil error", err)
			}

			if !tt.wantToPass && err == nil {
				t.Errorf("CreateTableMetadata() Error: nil, want error")
			}

		})
	}
}
