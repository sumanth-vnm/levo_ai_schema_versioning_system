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
	fileWriter, err := writer.CreateFormFile("file", "dummy.json")
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
	fileStore := storage.NewFileStore("test_dummy_directory")

	apiHandler := NewAPIHandler(fileStore, database)

	// Call the handler function
	fmt.Println(req)
	apiHandler.UploadSchemaHandler(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	expectedResponse := `{"message":"Schema uploaded successfully","version":3}` // CHANGE THIS EVERY TIME SINCE VERSION INCREMENTS
	if rr.Body.String() != expectedResponse {
		t.Errorf("expected response '%s' but got '%s'", expectedResponse, rr.Body.String())
	}
}


func TestGetLatestSchemaHandler(t *testing.T) {
	// Create a mock HTTP request with path variables
	req, err := http.NewRequest("GET", "/getLatestSchema/dummy.json", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set path variables in the request's context
	req = mux.SetURLVars(req, map[string]string{
		"filename": "dummy.json",
	})

	// Create a mock HTTP response recorder
	rr := httptest.NewRecorder()

	// Initialize the database connection
	database, err := db.Initialize()
	if err != nil {
		fmt.Errorf("Failed to initialize database: %v", err)
	}
	defer database.DB.Close()

	fileStore := storage.NewFileStore("test_dummy_directory")

	apiHandler := NewAPIHandler(fileStore, database)

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

	if _, ok := resp["file-dummy.json"]; !ok {
		t.Error("expected 'file-dummy.json' key in the response")
	}
}


func TestGetAllVersionsHandler(t *testing.T) {
	// Create a mock HTTP request with path variables
	req, err := http.NewRequest("GET", "/getAllVersions/dummy.json", nil)
	if err != nil {
		t.Fatal(err)
	}

	req = mux.SetURLVars(req, map[string]string{
		"filename": "dummy.json",
	})

	// Create a mock HTTP response recorder
	rr := httptest.NewRecorder()

	database, err := db.Initialize()
	if err != nil {
		fmt.Errorf("Failed to initialize database: %v", err)
	}
	defer database.DB.Close()
	
	fileStore := storage.NewFileStore("test_dummy_directory")

	apiHandler := NewAPIHandler(fileStore, database)

	// Call the handler function
	apiHandler.GetAllVersionsHandler(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, rr.Code)
	}


	// Parse the response JSON
	fmt.Println()
	fmt.Println("response body", rr.Body.String())

	for _, b := range rr.Body.Bytes() {
		if b < 0 || b > 127 {
			t.Fatalf("Cannot convert []byte to []int64: byte value out of range")
			return
		}
	}

	var versions map[string]interface{}
	var intf interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &intf)
	if err != nil {
		t.Fatalf("failed to unmarshal response JSON: %v", err)
	}
	versions = intf.(map[string]interface{})
	fmt.Println("versions", versions)
}

