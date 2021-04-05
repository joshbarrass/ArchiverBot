package internal

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// TestWritePermission checks whether write permission is available in
// a directory. This is done by writing to a new file in this directory.
func TestWritePermission(dir string) bool {
	filename := filepath.Join(dir, uuid.New().String())
	emptyFile, err := os.Create(filename)
	if err != nil {
		return false
	}
	emptyFile.Close()
	err = os.Remove(filename)
	if err != nil {
		return false
	}
	return true
}

// TestReadPermission checks whether read permission is available in a
// directory
func TestReadPermission(dir string) bool {
	_, err := ioutil.ReadDir(dir)
	if err != nil {
		return false
	}
	return true
}
