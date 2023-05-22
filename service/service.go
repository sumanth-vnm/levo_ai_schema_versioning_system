package service

import (
	"encoding/json"
	"fmt"
	yaml "gopkg.in/yaml.v2"

)

// ValidateJSONSchema validates the JSON schema file
func ValidateJSONSchema(schemaFile []byte) error {
	var data interface{}
	err := json.Unmarshal(schemaFile, &data)
	if err != nil {
		fmt.Println("failed to parse JSON schema:", err)
		return fmt.Errorf("failed to parse JSON schema: %v", err)
	}

	return nil
}

// ValidateYAMLSchema validates the YAML schema file
func ValidateYAMLSchema(schemaFile []byte) error {
	var data interface{}
	err := yaml.Unmarshal(schemaFile, &data)
	if err != nil {
		fmt.Println("failed to parse YAML schema:", err)
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
		fmt.Println("unsupported file type:", fileType)
		return fmt.Errorf("unsupported file type: %s", fileType)
	}

	return nil
}
