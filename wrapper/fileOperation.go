package wrapper

import (
	"fmt"
	"os"
)

type FileOperator interface {
	Stat(name string) (os.FileInfo, error)
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
	WriteString(file *os.File, s string) (int, error)
	CloseFile(file *os.File) error
	DeleteFile(fileName string) error
}

// FileWrapper is the  implementation of FileOperator
// using the actual os package
type fileOperator struct {
}

// Stat returns the file information for the given file name using the os.Stat function.
// It takes one parameter:
// - name: the name of the file to retrieve information for.
//
// It returns the file information as an os.FileInfo structure and an error, if any occurred.
func (f *fileOperator) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// OpenFile opens a file by calling os.OpenFile function.
// It takes three parameters:
//   - name: the name of the file to be opened.
//   - flag: the flag used to open the file.
//   - perm: the permission mode of the file.
//
// It returns a pointer to the opened file and an error, if any occurred.
func (f *fileOperator) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

// CloseFile closes the specified file by calling the Close method on the file object.
// It takes one parameter:
//   - file: the file object to be closed.
//
// It returns an error, if any occurred.
func (f *fileOperator) CloseFile(file *os.File) error {
	return file.Close()
}

// WriteString writes to the specified file by calling fmt.Fprintf method
// It takes two parameter
// - file which is the file object where to write the string
// - s the string need to write to the file
// Returns
// - int number of bytes written
// - error if there is any or nil
func (f *fileOperator) WriteString(file *os.File, s string) (int, error) {
	return fmt.Fprintf(file, "%s\n", s)
}

func (f *fileOperator) DeleteFile(fileName string) error {
	return os.Remove(fileName)
}

func NewFileOperator() FileOperator {
	return &fileOperator{}
}
