package wrapper

import (
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

// MockFileOperator is a mock implementation of FileOperator
type MockFileOperator struct {
	mock.Mock
}

// Stat mocks the Stat method of FileOperator interface
func (m *MockFileOperator) Stat(name string) (os.FileInfo, error) {
	args := m.Called(name)
	return args.Get(0).(os.FileInfo), args.Error(1)
}

// OpenFile mocks the OpenFile method of FileOperator interface
func (m *MockFileOperator) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	args := m.Called(name, flag, perm)
	return args.Get(0).(*os.File), args.Error(1)
}

// CloseFile mocks the CloseFile method of FileOperator interface
func (m *MockFileOperator) CloseFile(file *os.File) error {
	args := m.Called(file)
	return args.Error(0)
}

// WriteString mocks the WriteString method of FileOperator interface
func (m *MockFileOperator) WriteString(file *os.File, s string) (int, error) {
	args := m.Called(file, s)
	return args.Int(0), args.Error(1)
}

func NewFileOperatorMock(t testing.TB) *MockFileOperator {
	mock := &MockFileOperator{}
	mock.Mock.Test(t)
	t.Cleanup(func() {
		mock.AssertExpectations(t)
	})
	return mock
}
