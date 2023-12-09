package database

import (
	"custom_db/database/mocks"
	"custom_db/wrapper"
	"errors"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

func Test_InsertIntoTable(t *testing.T) {
	type fields struct {
		fileOperator    *wrapper.MockFileOperator
		metadataHandler *mocks.MetadataHandler
	}

	type testCase struct {
		name       string
		fields     fields
		values     []string
		createErr  bool
		writeErr   bool
		beforeTest func(*fields)
		afterTest  func(*fields)
		wantErr    bool
	}

	tests := []testCase{
		{
			name:      "no error",
			fields:    fields{},
			values:    []string{"one", "two", "three"},
			createErr: false,
			beforeTest: func(f *fields) {
				var file os.File
				f.fileOperator.On(
					"OpenFile",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(&file, nil)

				f.metadataHandler.On(
					"ReadColumnTypes",
					mock.Anything,
				).Return(make(map[string]string), nil)

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
			afterTest: func(*fields) {
			},
			writeErr: false,
			wantErr:  false,
		},
		{
			name:   "open error",
			values: []string{"stooges", "bar"},
			beforeTest: func(f *fields) {
				var file *os.File
				f.fileOperator.On(
					"OpenFile",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(file, errors.New("error opening table file "))
			},
			afterTest: func(fields *fields) {

			},
			createErr: true,
			writeErr:  false,
			wantErr:   true,
		},
		{
			name:      "write error",
			values:    []string{"foo", "invalid"},
			createErr: false,
			beforeTest: func(f *fields) {
				var file os.File
				f.fileOperator.On(
					"OpenFile",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(&file, nil)

				f.metadataHandler.On(
					"ReadColumnTypes",
					mock.Anything,
				).Return(make(map[string]string), nil)

				f.fileOperator.On(
					"WriteString",
					mock.Anything,
					mock.Anything,
				).Return(int(0), errors.New("error writing to table file"))

				f.fileOperator.On(
					"CloseFile",
					mock.Anything,
				).Return(nil)
			},
			afterTest: func(*fields) {
			},
			writeErr: true,
			wantErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.fields.fileOperator = wrapper.NewFileOperatorMock(t)
			tc.fields.metadataHandler = mocks.NewMetadataHandler(t)

			if tc.beforeTest != nil {
				tc.beforeTest(&tc.fields)
			}
			handler := NewTableHandler(tc.fields.fileOperator, tc.fields.metadataHandler)

			gotErr := handler.InsertIntoTable(tc.values) != nil
			if gotErr != tc.wantErr {
				t.Errorf("InsertIntoTable() error = %v, wantErr %v", gotErr, tc.wantErr)
				return
			}
		})
	}
}
