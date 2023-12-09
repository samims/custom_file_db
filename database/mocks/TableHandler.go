// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// TableHandler is an autogenerated mock type for the TableHandler type
type TableHandler struct {
	mock.Mock
}

// InsertIntoTable provides a mock function with given fields: values
func (_m *TableHandler) InsertIntoTable(values []string) error {
	ret := _m.Called(values)

	if len(ret) == 0 {
		panic("no return value specified for InsertIntoTable")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]string) error); ok {
		r0 = rf(values)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateDataType provides a mock function with given fields: values
func (_m *TableHandler) ValidateDataType(values []string) error {
	ret := _m.Called(values)

	if len(ret) == 0 {
		panic("no return value specified for ValidateDataType")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]string) error); ok {
		r0 = rf(values)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTableHandler creates a new instance of TableHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTableHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *TableHandler {
	mock := &TableHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
