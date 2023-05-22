
package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"path"
)

// FileStore represents the file storage
type FileStore struct {
	BasePath string
}

// NewFileStore creates a new file store
func NewFileStore(basePath string) *FileStore {
	return &FileStore{BasePath: basePath}
}

func (fs *FileStore) SaveSchema(schemaFile []byte, filename string, filetype string, version int64) error {
	// Create a new directory for each new file
	dirPath := filepath.Join(fs.BasePath, filename)
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Generate a unique filename for each version of the file
	newFilename := strconv.FormatInt(version, 10) + "." + filetype

	filePath := filepath.Join(dirPath, newFilename)

	err = ioutil.WriteFile(filePath, schemaFile, 0644)
	if err != nil {
		return fmt.Errorf("failed to save schema file: %v", err)
	}

	return nil
}

// GetSchema retrieves the schema file from the file store
func (fs *FileStore) GetSchema(filename string, version int64) ([]byte, error) {
	dirPath := filepath.Join(fs.BasePath, filename)
	fileType := strings.ToLower(path.Ext(filename))
	filePath := filepath.Join(dirPath, strconv.FormatInt(version, 10) + fileType)

	schemaFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("schema file '%s' version '%d' does not exist", filename, version)
		}
		return nil, fmt.Errorf("failed to read schema file: %v", err)
	}

	return schemaFile, nil
}


// DeleteSchema deletes the schema file with the specified filename and version from the storage
func (fs *FileStore) DeleteSchema(filename string, version int64) error {
	fileType := strings.ToLower(path.Ext(filename))
	dirPath := filepath.Join(fs.BasePath, filename, strconv.FormatInt(version, 10) + fileType)
	err := os.RemoveAll(dirPath)
	if err != nil {
		return fmt.Errorf("failed to delete schema: %v", err)
	}

	fmt.Println("Schema deleted successfully from storage")

	return nil
}

