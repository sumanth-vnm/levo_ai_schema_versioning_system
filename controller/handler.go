package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"example.com/levo_app/db"
	"example.com/levo_app/storage"
	"example.com/levo_app/service"

	"github.com/gorilla/mux"
)

// APIHandler represents the API handler
type APIHandler struct {
	Storage  *storage.FileStore
	Database *db.Database
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(storage *storage.FileStore, database *db.Database) *APIHandler {
	return &APIHandler{
		Storage:  storage,
		Database: database,
	}
}

// UploadSchemaHandler handles the API for uploading a schema
func (ah *APIHandler) UploadSchemaHandler(w http.ResponseWriter, r *http.Request) {
	
	
	file, fileHeaders, err := r.FormFile("file")
	if err != nil {
		fmt.Println("failed to read file", err)
		http.Error(w, "failed to read file or 'file' field doesn't exist in request body", http.StatusBadRequest)
		return
	}
	defer file.Close()

	

	schemaFile, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "failed to read file 2", http.StatusInternalServerError)
		return
	}

	fmt.Println("File Read Succesfully")

	// filename := r.FormValue("filename")
	filename := fileHeaders.Filename
	fileType := strings.ToLower(path.Ext(filename))
	fileType = strings.TrimPrefix(fileType, ".")

	fmt.Println("File Type: ", fileType)
	fmt.Println("File name: ", filename)


	// Validate the schema file
	err = service.ValidateSchema(schemaFile, fileType)
	if err != nil {
		http.Error(w, fmt.Sprintf("INVALID SCHEMA: %v", err), http.StatusBadRequest)
		return
	}

	fmt.Println("File Validated Succesfully")

	// Get the latest version number from the database
	latestVersion, err := ah.Database.GetLatestSchemaVersion(filename)
	if err != nil {
		latestVersion = 0
	}

	fmt.Println("Latest version fetched Succesfully")
	version := latestVersion + 1

	err = ah.Storage.SaveSchema(schemaFile, filename, fileType, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Schema saved succesfully in storage")

	timestamp := time.Now()
	schema := db.Schema{
		Version:   version,
		Filename:  filename,
		Timestamp: timestamp,
	}

	err = ah.Database.SaveSchema(schema)
	if err != nil {
		http.Error(w, "failed to save schema", http.StatusInternalServerError)
		// remove from storage as well
		err := ah.Storage.DeleteSchema(schema.Filename, schema.Version)
		if err != nil {
			fmt.Println("failed to delete schema from storage:", err)
		}
		return
	}

	fmt.Println("Schema saved succesfully in database")

	// Give success response
	resp := make(map[string]interface{})
	resp["message"] = "Schema uploaded successfully"
	resp["version"] = schema.Version

	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "failed to marshal response to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}

// GetSchemaHandler handles the API for retrieving a specific version of a schema
func (ah *APIHandler) GetSchemaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	filename, ok := vars["filename"]
	if !ok {
		http.Error(w, "filename not found in request", http.StatusBadRequest)
		return
	}

	version, ok := vars["version"]
	if !ok {
		http.Error(w, "version not found in request", http.StatusBadRequest)
		return
	}
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, "Version is not an integer", http.StatusBadRequest)
		return
	}

	schema, err := ah.Database.GetSchema(filename, int64(versionInt))
	if err != nil {
		http.Error(w, "schema not found", http.StatusNotFound)
		return
	}

	schemaFile, err := ah.Storage.GetSchema(schema.Filename, schema.Version)
	if err != nil {
		http.Error(w, "failed to read schema file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(schemaFile)
}

func (ah *APIHandler) GetLatestSchemaHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Get Latest Schema Handler Called")
	vars := mux.Vars(r)
	filename, ok := vars["filename"]
	if !ok {
		http.Error(w, "filename not found in request", http.StatusBadRequest)
		return
	}
	fmt.Println("filename: ", filename)

	latestVersion, err := ah.Database.GetLatestSchemaVersion(filename)
	if err != nil {
		http.Error(w, "failed to get latest schema version", http.StatusInternalServerError)
		return
	}

	schemaFile, err := ah.Storage.GetSchema(filename, latestVersion)
	if err != nil {
		http.Error(w, "failed to read schema file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	// Creating response object
	resp := make(map[string]interface{})
	resp["version"] = latestVersion
	schemaFileJson := make(map[string]interface{})
	err = json.Unmarshal(schemaFile, &schemaFileJson)
	if err != nil {
		http.Error(w, "failed to ummarshal file to JSON", http.StatusInternalServerError)
		return
	}
	resp["file-" + filename] = schemaFileJson

	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "failed to marshal response to JSON", http.StatusInternalServerError)
		return
	}

	w.Write(respBytes)
	
}

func (ah *APIHandler) GetAllVersionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename, ok := vars["filename"]
	if !ok {
		http.Error(w, "filename not found in request", http.StatusBadRequest)
		return
	}


	fmt.Println("filename: ", filename)

	// Call the storage method to retrieve the versions for the specified filename
	versions, err := ah.Database.GetAllVersionsForSchema(filename)
	if err != nil {
		http.Error(w, "failed to get versions for schema", http.StatusInternalServerError)
		return
	}

	resp := make(map[string]interface{})
	resp["available_versions"] = versions

	// Convert the versions to JSON
	jsonVersions, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "failed to marshal versions to JSON", http.StatusInternalServerError)
		return
	}

	

	// Set the response headers and write the JSON versions
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonVersions)
}


