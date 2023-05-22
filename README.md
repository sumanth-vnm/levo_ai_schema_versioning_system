# levo_ai_schema_versioning_system

An Upload/Version Schema API implementation using Go lang.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)

## Installation

### GO project initialization

Follow these steps to clone and set up the project:

1. Clone the GitHub repository using the following command:

   ```terminal
   gh repo clone sumanth-vnm/levo_ai_schema_versioning_system
   ```

2. Open a terminal in the installed repos' root directory and initialize go module using

    ```terminal
   go mod init example.com/levo_app
   ```

3. Install dependencies using

    ```terminal
   go mod tidy
   ```


### Postgres database connection

Follow the below steps, to set up the postgres database connection to the go project:

1. Use any existing database or create a new database in postgres. Let's say we are using a database named "levo_app_db"

2. create a new sequence in the database using the following command:

   ```terminal
   CREATE SEQUENCE levo_sequence
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
   ```

3. create a new table in the database using the following command (This table uses the above sequence to generate the primary key):

   ```terminal
    CREATE TABLE schemas(
    id BIGINT PRIMARY KEY DEFAULT nextval('levo_sequence'),
    filename TEXT,
    version BIGINT,
    created_on TIMESTAMPTZ
    );
    ```

4. Inside the cloned project, in the file "db/database.go",
    update the database connection parameters which looks like below:
    
    ```terminal
    username := "postgres"
	password := "sam123"
	dbName := "levo_app_db"
    ```

## Usage

### Run the project

1. to run this project, use the following command:

    ```terminal
    go run main.go
    ```
2. This will start the server at port 8080. if port needs to be changed, update the port number in the file "main.go" which looks like below:

    ```terminal
    http.ListenAndServe(":8080", router)
    ```

### API endpoints

Following are the API endpoints implemented in this project:

```terminal
localhost:8080/upload/schema - (POST) - to upload a new schema file[json/yaml]
    Note: Use the field key "file" inside body to upload any schema file
    Output: If successful, returns the "version" and success "message"

localhost:8080/getLatestSchema/{{filename}} - (GET) - to get the latest schema file
     Output: If successful, returns the latest schema file with field "filename" and "version" number

localhost:8080/getSchemaByVersion/{{filename}}/{version} - (GET) - to get the schema file of a particular file name and version
     Output: If successful, returns the schema file of requested version number

localhost:8080/getAllVersions/{{filename}} - (GET) - to get all the schema files
     Output: If successful, returns "available_versions" array with all the versions of the requested file name
```

1. Once a schema is uploaded, The uploaded files will be stored under "schema_uploads" folder in the root directory of the project.

2. For every new schema file uploaded, a new directory will be created with the name of the file under the "schema_uploads" folder.

3. Inside the directory of the uploaded file, naming of the different versions will be 1.json, 2.json, 3.json, etc. or with yaml type extensions.

### Postman collection

a postman collection json file is added in the repository at the root folder named "Levo.ai.postman_collection.json". This file can be imported into postman to test the API endpoints.





 
