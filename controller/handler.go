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
	yaml "gopkg.in/yaml.v2"

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

// ValidateJSONSchema validates the JSON schema file
func ValidateJSONSchema(schemaFile []byte) error {
	var data interface{}
	err := json.Unmarshal(schemaFile, &data)
	if err != nil {
		return fmt.Errorf("failed to parse JSON schema: %v", err)
	}

	return nil
}

// ValidateYAMLSchema validates the YAML schema file
func ValidateYAMLSchema(schemaFile []byte) error {
	var data interface{}
	err := yaml.Unmarshal(schemaFile, &data)
	if err != nil {
		return fmt.Errorf("failed to parse YAML schema: %v", err)
	}

	return nil
}

// ValidateSchema validates the schema file based on its type (JSON or YAML)
func ValidateSchema(schemaFile []byte, fileType string) error {
	if fileType == "json" {
		// Validate JSON schema
		err := ValidateJSONSchema(schemaFile)
		if err != nil {
			return err
		}
	} else if fileType == "yaml" {
		// Validate YAML schema
		err := ValidateYAMLSchema(schemaFile)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported file type: %s", fileType)
	}

	return nil
}

// UploadSchemaHandler handles the API for uploading a schema
func (ah *APIHandler) UploadSchemaHandler(w http.ResponseWriter, r *http.Request) {
	
	fmt.Println("Upload Schema Handler Called")
	file, fileHeaders, err := r.FormFile("file")
	fmt.Println("FormFile Called")
	if err != nil {
		fmt.Println("failed to read file", err)
		http.Error(w, "failed to read file 1", http.StatusBadRequest)
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
	err = ValidateSchema(schemaFile, fileType)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid schema: %v", err), http.StatusBadRequest)
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

	err = ah.Storage.SaveSchema(schemaFile, filename, version)
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
		// TODO remove from storage
		return
	}

	fmt.Println("Schema saved succesfully in database")

	fmt.Fprintf(w, "Schema uploaded successfully")
	// TODO give success response
}

// GetSchemaHandler handles the API for retrieving a specific version of a schema
func (ah *APIHandler) GetSchemaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]
	version := vars["version"]
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
	filename := vars["filename"] // TODO error handling
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
	filename := vars["filename"]


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


