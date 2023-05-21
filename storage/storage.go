
package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// FileStore represents the file storage
type FileStore struct {
	BasePath string
}

// NewFileStore creates a new file store
func NewFileStore(basePath string) *FileStore {
	return &FileStore{BasePath: basePath}
}

func (fs *FileStore) SaveSchema(schemaFile []byte, filename string, version int64) error {
	// Create a new directory for each new file
	dirPath := filepath.Join(fs.BasePath, filename)
	err := os.MkdirAll(dirPath, 0755) // TODO check permissions
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Generate a unique filename for each version of the file
	newFilename := strconv.FormatInt(version, 10)

	filePath := filepath.Join(dirPath, newFilename)

	err = ioutil.WriteFile(filePath, schemaFile, 0644) // TODO check permissions
	if err != nil {
		return fmt.Errorf("failed to save schema file: %v", err)
	}

	return nil
}

// GetSchema retrieves the schema file from the file store
func (fs *FileStore) GetSchema(filename string, version int64) ([]byte, error) {
	dirPath := filepath.Join(fs.BasePath, filename)
	filePath := filepath.Join(dirPath, strconv.FormatInt(version, 10))

	schemaFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("schema file '%s' version '%d' does not exist", filename, version)
		}
		return nil, fmt.Errorf("failed to read schema file: %v", err)
	}

	return schemaFile, nil
}



// GetLatestSchema retrieves the latest schema file from the file store
func (fs *FileStore) GetLatestSchema(filename string) ([]byte, error) {
	dirPath := filepath.Join(fs.BasePath, filename)
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema files: %v", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no schema files found")
	}

	latestSchemaFile := files[len(files)-1]
	latestSchemaFilename := latestSchemaFile.Name()

	filePath := filepath.Join(dirPath, latestSchemaFilename)

	schemaFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %v", err)
	}

	return schemaFile, nil
}
