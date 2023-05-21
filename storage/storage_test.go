package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveSchema(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new instance of FileStore with the temporary directory as the base path
	fileStore := &FileStore{
		BasePath: tempDir,
	}

	// Define test data
	schemaFile := []byte("test schema")
	filename := "openapi.json"
	version := int64(1)

	// Call the SaveSchema function
	err = fileStore.SaveSchema(schemaFile, filename, version)
	if err != nil {
		t.Errorf("failed to save schema: %v", err)
	}

	// Verify that the directory and file were created
	dirPath := filepath.Join(tempDir, filename)
	filePath := filepath.Join(dirPath, "1")
	if _, err := os.Stat(dirPath); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("expected directory '%s' to be created but it does not exist", dirPath)
		} else {
			t.Errorf("failed to access directory '%s': %v", dirPath, err)
		}
	}
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("expected file '%s' to be created but it does not exist", filePath)
		} else {
			t.Errorf("failed to access file '%s': %v", filePath, err)
		}
	}

	// Verify the content of the file
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file '%s': %v", filePath, err)
	}
	if string(fileContent) != string(schemaFile) {
		t.Errorf("expected file content to be '%s' but got '%s'", string(schemaFile), string(fileContent))
	}
}

// Rest of the struct definitions and function implementations remain the same as before

func TestGetSchema(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new instance of FileStore with the temporary directory as the base path
	fileStore := &FileStore{
		BasePath: tempDir,
	}

	// Create a test schema file
	schemaFileContent := []byte("test schema")
	filename := "openapi.json"
	version := int64(1)
	filePath := filepath.Join(tempDir, filename, "1")
	err = os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(filePath, schemaFileContent, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Call the GetSchema function
	schemaFile, err := fileStore.GetSchema(filename, version)
	if err != nil {
		t.Errorf("failed to get schema: %v", err)
	}

	// Verify the content of the schema file
	if string(schemaFile) != string(schemaFileContent) {
		t.Errorf("expected schema file content to be '%s' but got '%s'", string(schemaFileContent), string(schemaFile))
	}

	// Call the GetSchema function with a non-existing version
	nonExistingVersion := int64(2)
	_, err = fileStore.GetSchema(filename, nonExistingVersion)
	if err == nil {
		t.Errorf("expected error when getting non-existing schema file version, but got no error")
	} else {
		expectedError := fmt.Sprintf("schema file '%s' version '%d' does not exist", filename, nonExistingVersion)
		if err.Error() != expectedError {
			t.Errorf("expected error message '%s' but got '%s'", expectedError, err.Error())
		}
	}
}


func TestGetLatestSchema(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new instance of FileStore with the temporary directory as the base path
	fileStore := &FileStore{
		BasePath: tempDir,
	}

	// Create multiple test schema files
	filename := "openapi.json"
	schemaFileContent1 := []byte("test schema 1")
	schemaFileContent2 := []byte("test schema 2")
	schemaFileContent3 := []byte("test schema 3")

	// Save the test schema files in the temporary directory
	filePath1 := filepath.Join(tempDir, filename, "1")
	filePath2 := filepath.Join(tempDir, filename, "2")
	filePath3 := filepath.Join(tempDir, filename, "3")
	err = os.MkdirAll(filepath.Dir(filePath1), 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(filepath.Dir(filePath2), 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(filepath.Dir(filePath3), 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(filePath1, schemaFileContent1, 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(filePath2, schemaFileContent2, 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(filePath3, schemaFileContent3, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Call the GetLatestSchema function
	schemaFile, err := fileStore.GetLatestSchema(filename)
	if err != nil {
		t.Errorf("failed to get latest schema: %v", err)
	}

	// Verify the content of the latest schema file
	if string(schemaFile) != string(schemaFileContent3) {
		t.Errorf("expected latest schema file content to be '%s' but got '%s'", string(schemaFileContent3), string(schemaFile))
	}

	// Call the GetLatestSchema function when no schema files exist
	emptyDir := filepath.Join(tempDir, "empty")
	err = os.MkdirAll(emptyDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	_, err = fileStore.GetLatestSchema("empty")
	if err == nil {
		t.Errorf("expected error when getting latest schema from empty directory, but got no error")
	} else {
		expectedError := "no schema files found"
		if err.Error() != expectedError {
			t.Errorf("expected error message '%s' but got '%s'", expectedError, err.Error())
		}
	}
}
