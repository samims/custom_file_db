// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MetadataHandler is an autogenerated mock type for the MetadataHandler type
type MetadataHandler struct {
	mock.Mock
}

func (_m *MetadataHandler) ReadColNamesAndTypesInArray(fileName string) ([]string, []string, error) {
	//TODO implement me
	panic("implement me")
}

// CreateTableMetadata provides a mock function with given fields: colNames, colTypes
func (_m *MetadataHandler) CreateTableMetadata(colNames []string, colTypes []string) error {
	ret := _m.Called(colNames, colTypes)

	if len(ret) == 0 {
		panic("no return value specified for CreateTableMetadata")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]string, []string) error); ok {
		r0 = rf(colNames, colTypes)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReadColumnTypes provides a mock function with given fields: filename
func (_m *MetadataHandler) ReadColumnTypes(filename string) (map[string]string, error) {
	ret := _m.Called(filename)

	if len(ret) == 0 {
		panic("no return value specified for ReadColumnTypes")
	}

	var r0 map[string]string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (map[string]string, error)); ok {
		return rf(filename)
	}
	if rf, ok := ret.Get(0).(func(string) map[string]string); ok {
		r0 = rf(filename)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(filename)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMetadataHandler creates a new instance of MetadataHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMetadataHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *MetadataHandler {
	mock := &MetadataHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
