package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Database represents the database
type Database struct {
	DB *sql.DB
}

// Schema represents the schema record in the database
type Schema struct {
	ID        int
	Version   int64
	Filename  string
	Timestamp time.Time
}

// Initialize initializes the database connection
func Initialize() (*Database, error) {
	username := "postgres"
	password := "sam123"
	dbName := "postgres"
	connStr := "postgres://" + username + ":" + password + "@localhost/" + dbName + "?sslmode=disable" // TODO - Update the connection string
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping the database: %v", err)
	}

	log.Println("Connected to the database")

	return &Database{DB: db}, nil
}

// SaveSchema saves the schema record to the database
func (db *Database) SaveSchema(schema Schema) error {
	fmt.Println("Saving schema...")
	fmt.Println("schema details", schema.Version, schema.Filename, schema.Timestamp)
	query := "INSERT INTO schemas (version, filename, created_on) VALUES ($1, $2, $3)"
	_, err := db.DB.Exec(query, schema.Version, schema.Filename, schema.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to save schema: %v", err)
	}

	return nil
}

// GetSchema retrieves a specific version of a schema from the database
func (db *Database) GetSchema(filename string, version int64) (Schema, error) {
	query := "SELECT id, version, filename, created_on FROM schemas WHERE filename = $1 AND version = $2"
	row := db.DB.QueryRow(query, filename, version)

	var schema Schema
	err := row.Scan(&schema.ID, &schema.Version, &schema.Filename, &schema.Timestamp)
	if err != nil {
		return Schema{}, fmt.Errorf("failed to get schema: %v", err)
	}

	return schema, nil
}

func (db *Database) GetLatestSchemaVersion(filename string) (int64, error) {
	fmt.Println("Getting latest schema version...")
	query := "SELECT MAX(version) FROM schemas WHERE filename = $1"
	row := db.DB.QueryRow(query, filename)

	var latestVersion sql.NullInt64
	err := row.Scan(&latestVersion)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest schema version: %v", err)
	}

	fmt.Println("latestVersion: ", latestVersion)

	if latestVersion.Valid {
		return latestVersion.Int64, nil
	}

	return 0, nil
}

// GetAllVersionsForSchema retrieves all available versions for a specific schema filename from the database
func (db *Database) GetAllVersionsForSchema(filename string) ([]int64, error) {
	// Execute a query to retrieve the versions for the given filename from the database
	// Here's an example using PostgreSQL as the database

	// Assuming you have a table named 'schema_versions' with columns named 'filename' and 'version'
	// and you want to retrieve all versions for a specific filename
	query := "SELECT version FROM schemas WHERE filename = $1"

	var versions []int64
	rows, err := db.DB.Query(query, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve versions for file '%s': %v", filename, err)
	}
	defer rows.Close()

	fmt.Println("rows: ", rows)

	for rows.Next() {
		var version int64
		err := rows.Scan(&version)
		if err != nil {
			return nil, fmt.Errorf("failed to scan version: %v", err)
		}
		versions = append(versions, version)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over versions: %v", err)
	}

	return versions, nil
}
