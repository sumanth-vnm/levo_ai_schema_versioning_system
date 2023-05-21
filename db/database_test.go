package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"


	_ "github.com/lib/pq" // PostgreSQL driver
)

func TestInitialize(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database connection: %v", err)
	}
	defer db.Close()

	// Set the expectations for the mock database
	mock.ExpectPing()

	// Call the Initialize function with the mock database
	database, err := InitializeWithDB(db)
	if err != nil {
		t.Errorf("unexpected error during initialization: %v", err)
	}

	// Verify that the database connection is set
	if database.DB != db {
		t.Errorf("expected database connection to be set, but it was not")
	}

	// Verify that the expected methods were called on the mock
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func InitializeWithDB(db *sql.DB) (*Database, error) {
	err := db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping the database: %v", err)
	}

	log.Println("Connected to the database")

	return &Database{DB: db}, nil
}

func TestSaveSchema(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database connection: %v", err)
	}
	defer db.Close()

	// Create the schema object for testing
	schema := Schema{
		Version:   1,
		Filename:  "test_schema.json",
		Timestamp: time.Now(),
	}

	// Set the expectations for the mock database
	mock.ExpectExec("INSERT INTO schemas").
		WithArgs(schema.Version, schema.Filename, schema.Timestamp).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create a database instance with the mock connection
	database := &Database{DB: db}

	// Call the SaveSchema function with the schema object
	err = database.SaveSchema(schema)
	if err != nil {
		t.Errorf("unexpected error while saving schema: %v", err)
	}

	// Verify that the expected query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}


func TestGetSchema(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database connection: %v", err)
	}
	defer db.Close()

	// Create the expected schema object for testing
	expectedSchema := Schema{
		ID:        1,
		Version:   1,
		Filename:  "test_schema.json",
		Timestamp: time.Now(),
	}

	// Set the expectations for the mock database
	mock.ExpectQuery("SELECT id, version, filename, created_on FROM schemas").
		WithArgs(expectedSchema.Filename, expectedSchema.Version).
		WillReturnRows(sqlmock.NewRows([]string{"id", "version", "filename", "created_on"}).
			AddRow(expectedSchema.ID, expectedSchema.Version, expectedSchema.Filename, expectedSchema.Timestamp))

	// Create a database instance with the mock connection
	database := &Database{DB: db}

	// Call the GetSchema function with the filename and version
	schema, err := database.GetSchema(expectedSchema.Filename, expectedSchema.Version)
	if err != nil {
		t.Errorf("unexpected error while getting schema: %v", err)
	}

	// Verify that the returned schema matches the expected schema
	if schema != expectedSchema {
		t.Errorf("unexpected schema returned, expected: %+v, got: %+v", expectedSchema, schema)
	}

	// Verify that the expected query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}


func TestGetLatestSchemaVersion(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database connection: %v", err)
	}
	defer db.Close()

	// Create the expected schema object for testing
	expectedSchema := Schema{
		ID:        1,
		Version:   1,
		Filename:  "test_schema.json",
		Timestamp: time.Now(),
	}

	// Set the expectations for the mock database
	mock.ExpectQuery("SELECT id, version, filename, created_on FROM schemas").
		WithArgs(expectedSchema.Filename, expectedSchema.Version).
		WillReturnRows(sqlmock.NewRows([]string{"id", "version", "filename", "created_on"}).
			AddRow(expectedSchema.ID, expectedSchema.Version, expectedSchema.Filename, expectedSchema.Timestamp))

	// Set the expected filename and latest version
	expectedFilename := "test_schema.sql"
	expectedVersion := int64(1)

	// Set the expectations for the mock database
	mock.ExpectQuery("SELECT MAX(version) FROM schemas WHERE filename = $1").
		WithArgs(expectedFilename).
		WillReturnRows(sqlmock.NewRows([]string{"max"}).AddRow(expectedVersion))

	// Create a database instance with the mock connection
	database := &Database{DB: db}

	// Call the GetLatestSchemaVersion function with the filename
	latestVersion, err := database.GetLatestSchemaVersion(expectedFilename)
	fmt.Println("~~~", latestVersion)
	if err != nil {
		t.Errorf("unexpected error while getting latest schema version: %v", err)
	}

	// Verify that the returned latest version matches the expected version
	if latestVersion != expectedVersion {
		t.Errorf("unexpected latest schema version returned, expected: %d, got: %d", expectedVersion, latestVersion)
	}

	// Verify that the expected query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}


func TestGetAllVersionsForSchema(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database connection: %v", err)
	}
	defer db.Close()

	// Set the expected filename and versions
	expectedFilename := "test_schema.sql"
	expectedVersions := []int64{1, 2, 3}

	// Set the expectations for the mock database
	mock.ExpectQuery("SELECT version FROM schemas WHERE filename = $1").
		WithArgs(expectedFilename).
		WillReturnRows(sqlmock.NewRows([]string{"version"}).
			AddRow(expectedVersions[0]).
			AddRow(expectedVersions[1]).
			AddRow(expectedVersions[2]))

	// Create a database instance with the mock connection
	database := &Database{DB: db}

	// Call the GetAllVersionsForSchema function with the filename
	versions, err := database.GetAllVersionsForSchema(expectedFilename)
	if err != nil {
		t.Errorf("unexpected error while getting versions for schema: %v", err)
	}

	// Verify that the returned versions match the expected versions
	if len(versions) != len(expectedVersions) {
		t.Errorf("unexpected number of versions returned, expected: %d, got: %d", len(expectedVersions), len(versions))
	} else {
		for i, v := range versions {
			if v != expectedVersions[i] {
				t.Errorf("unexpected version at index %d, expected: %d, got: %d", i, expectedVersions[i], v)
			}
		}
	}

	// Verify that the expected query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}