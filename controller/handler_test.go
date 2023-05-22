package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/levo_app/db"
	"example.com/levo_app/storage"
	"github.com/gorilla/mux"
)

func TestUploadSchemaHandler(t *testing.T) {
	// Create a mock HTTP request with a file
	// Create a buffer to store the request body
	bodyBuf := &bytes.Buffer{}

	// Create a new multipart writer
	writer := multipart.NewWriter(bodyBuf)

	// Create a file part with the desired content and filename
	fileContents := []byte(`{
		"openapi": "3.0.1",
		"info": {
		  "title": "OWASP crAPI API",
		  "version": "1.0"
		}
		}`)
	fileWriter, err := writer.CreateFormFile("file", "openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	_, err = fileWriter.Write(fileContents)
	if err != nil {
		t.Fatal(err)
	}

	// Close the multipart writer to finalize the request body
	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/upload/schema", bodyBuf)
	if err != nil {
		t.Fatal(err)
	}

	// req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a mock HTTP response recorder
	rr := httptest.NewRecorder()

	// Initialize the database connection
	database, err := db.Initialize()
	if err != nil {
		fmt.Errorf("Failed to initialize database: %v", err)
	}
	defer database.DB.Close()

	// Initialize the file storage
	fileStore := storage.NewFileStore("schema_files")

	apiHandler := NewAPIHandler(fileStore, database)

	// Call the handler function
	fmt.Println(req)
	apiHandler.UploadSchemaHandler(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	expectedResponse := "Schema uploaded successfully"
	if rr.Body.String() != expectedResponse {
		t.Errorf("expected response '%s' but got '%s'", expectedResponse, rr.Body.String())
	}
}

// Mock implementation of the Database interface
type mockDatabase struct{}

func (db mockDatabase) GetLatestSchemaVersion(filename string) (int, error) {
	// Return a dummy latest version
	return 1, nil
}

func (db mockDatabase) SaveSchema(schema db.Schema) error {
	// Do nothing in the mock implementation
	return nil
}

// Mock implementation of the Storage interface
type mockStorage struct{}

func (s mockStorage) SaveSchema(schemaFile []byte, filename string, version int) error {
	// Do nothing in the mock implementation
	return nil
}

func (s mockStorage) DeleteSchema(filename string, version int) error {
	// Do nothing in the mock implementation
	return nil
}


func TestGetSchemaHandler(t *testing.T) {
	// Create a mock HTTP request with path variables
	req, err := http.NewRequest("GET", "/schema/openapi.json/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set path variables in the request's context
	req = mux.SetURLVars(req, map[string]string{
		"filename": "openapi.json",
		"version":  "1",
	})

	// Create a mock HTTP response recorder
	rr := httptest.NewRecorder()

	// Create an instance of APIHandler with necessary dependencies
	apiHandler := &APIHandler{
		Database: mockDatabase{}, // Provide a mock implementation of the Database interface
		Storage:  mockStorage{},  // Provide a mock implementation of the Storage interface
	}

	// Call the handler function
	apiHandler.GetSchemaHandler(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body or other assertions as needed
	// ...
}

// Rest of the mock implementations and struct definitions remain the same as before

func TestGetLatestSchemaHandler(t *testing.T) {
	// Create a mock HTTP request with path variables
	req, err := http.NewRequest("GET", "/latest-schema/openapi.json", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set path variables in the request's context
	req = mux.SetURLVars(req, map[string]string{
		"filename": "openapi.json",
	})

	// Create a mock HTTP response recorder
	rr := httptest.NewRecorder()

	// Create an instance of APIHandler with necessary dependencies
	apiHandler := &APIHandler{
		Database: mockDatabase{}, // Provide a mock implementation of the Database interface
		Storage:  mockStorage{},  // Provide a mock implementation of the Storage interface
	}

	// Call the handler function
	apiHandler.GetLatestSchemaHandler(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body or other assertions as needed
	// ...

	// Parse the response JSON
	var resp map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("failed to unmarshal response JSON: %v", err)
	}

	// Verify the expected structure of the response
	if _, ok := resp["version"]; !ok {
		t.Error("expected 'version' key in the response")
	}

	if _, ok := resp["file-openapi.json"]; !ok {
		t.Error("expected 'file-openapi.json' key in the response")
	}

	// ...
}

// Rest of the mock implementations and struct definitions remain the same as before

func TestGetAllVersionsHandler(t *testing.T) {
	// Create a mock HTTP request with path variables
	req, err := http.NewRequest("GET", "/versions/openapi.json", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set path variables in the request's context
	req = mux.SetURLVars(req, map[string]string{
		"filename": "openapi.json",
	})

	// Create a mock HTTP response recorder
	rr := httptest.NewRecorder()

	// Create an instance of APIHandler with necessary dependencies
	apiHandler := &APIHandler{
		Database: mockDatabase{}, // Provide a mock implementation of the Database interface
	}

	// Call the handler function
	apiHandler.GetAllVersionsHandler(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body or other assertions as needed
	// ...

	// Parse the response JSON
	var versions []int64
	err = json.Unmarshal(rr.Body.Bytes(), &versions)
	if err != nil {
		t.Fatalf("failed to unmarshal response JSON: %v", err)
	}

	// Verify the expected structure of the response
	// ...

	// ...
}

// Rest of the mock implementations and struct definitions remain the same as before
