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

2. Open a terminal in the installed repos' root directory and initialize go module using "go mod init example.com/levo_app"

3. install dependencies using "go mod tidy"


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






 
