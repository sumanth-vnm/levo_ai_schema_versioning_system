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

localhost:8080/getLatestSchema/{{filename}} - (GET) - to get the latest schema file

localhost:8080/getSchemaByVersion/{{filename}}/{version} - (GET) - to get the schema file of a particular file name and version

localhost:8080/getAllVersions/{{filename}} - (GET) - to get all the schema files
```

### Postman collection

The following postman collection can be used to test the API endpoints:

```terminal

```





 
