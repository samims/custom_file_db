package database

import (
	"custom_db/constants"
	"custom_db/wrapper"
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

func TestCreateTableMetadata(t *testing.T) {

	type fields struct {
		fileOperator *wrapper.MockFileOperator
	}

	type args struct {
		colNames []string
		colTypes []string
	}

	cleanUp := func(f *fields) {
		err := os.Remove(constants.DefaultTableMetadataName + ".txt")
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
			fields{},
			args{
				colNames: []string{"id", "name"},
				colTypes: []string{"int", "string"},
			},
			true,
			func(f *fields) {
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
				[]string{"id", "name"},
				[]string{"int", "string"},
			},
			false,
			func(f *fields) {
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
				[]string{"id"},
				[]string{"int", "string"},
			},
			false,
			func(f *fields) {
				//fileInfo :=
				f.fileOperator.On(
					"Stat",
					mock.Anything,
				).Return(&wrapper.MockFileInfo{}, errors.New("test"))
			},
			func(f *fields) {},
		},
		{
			"ErrorOpeningFile",
			fields{},
			args{
				[]string{"id"},
				[]string{"int"},
			},
			false,
			func(f *fields) {
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
			func(f *fields) {
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.fields = fields{
				fileOperator: wrapper.NewFileOperatorMock(t),
			}

			if tt.beforeTest != nil {
				tt.beforeTest(&tt.fields)
			}
			defer tt.afterTest(&tt.fields)

			mh := NewMetadataHandler(tt.fields.fileOperator)
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
