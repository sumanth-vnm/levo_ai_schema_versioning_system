package main

import (
	"log"
	"net/http"

	"example.com/levo_app/api"
	"example.com/levo_app/db"
	"example.com/levo_app/controller"
	"example.com/levo_app/storage"
	
)

func main() {
	// Initialize the database connection
	database, err := db.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.DB.Close()

	// Initialize the file storage
	fileStore := storage.NewFileStore("schema_uploads") // Give the base path as param in NewFileStore

	// Create the API handler
	apiHandler := controller.NewAPIHandler(fileStore, database)

	// Register API routes
	router := api.RegisterRoutes(apiHandler)

	log.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
