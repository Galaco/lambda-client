package filesystem

import (
	"fmt"
)

// FileNotFoundError
type FileNotFoundError struct {
	fileName string
}

// Error
func (err FileNotFoundError) Error() string {
	return fmt.Sprintf("%s not found in filesystem", err.fileName)
}

// NewFileNotFoundError
func NewFileNotFoundError(filename string) *FileNotFoundError {
	return &FileNotFoundError{
		fileName: filename,
	}
}
