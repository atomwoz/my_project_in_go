# My project for Remitly!

This is contenarized application that provides a REST API to deal with SWIFT codes.

## Features

- Uses GO - new, modern and fast language
- Uses Postgres - a well known powerful and fast database
- Uses Docker - to make app running in any environment, and fast to deploy
- It has tests - to ensure the quality of the code

## How to run

### Requirements

- Docker installed, instructions [here](https://docs.docker.com/get-docker/). At Linux based systems please don't forget to add your user to the docker group.

### 1. Clone the repository

```shell
git clone https://github.com/atomwoz/my_project_in_go
cd my_project_in_go
```

### 2. Modify the .env file

For safety, you must change the default password for the database.  
It is also recomended to change the database user too.  
NOTE: Everything in `.env` could be overriden by environment variables.

```shell
vim .env
```

### 3. Import CSV to the database

It will import the CSV file to the database.  
The file is located at `./data/swift_data.csv`. If you want to import another file, just place path to it as the first argument.  
At the first run, it will build the image and setup database.

```shell
docker compose run csv_import
```

### 4. Run the application

It will listen at [http://localhost:8080](http://localhost:8080)

```shell
docker compose up
```

## Running Tests

To run tests, you need to have Go installed.  
This project was tested with Go version 1.24.  
You can find the installation instructions [here](https://golang.org/doc/install)

To run integration tests:  
NOTE: Connection to the production database is required to run integration tests.

```shell
go test ./tests/...
```

To run unit tests:

```shell
go test ./internal/...
```

To test the connection, you can use the following command:

```shell
docker compose run db_check
```

To run all tests:

```shell
go test ./...
```

## Error format

### GET requests

They return errors with the following structure:

```json
{ "error": "<code>", "error_message": "<string message>" }
```

Where code is an integer and message is a string.  
Possible error codes are:

1: Internal database error  
2: SWIFT record not found  
3: Country code not found

### POST and DELETE requests

It returns only messages with the following structure:

```json
{ "message": "<string message>" }
```

Where message is a string.

If the request is successful, a `POST` will return a `201` status code with the following message:

```json
{ "message": "ok" }
```

Similarly, a successful `DELETE` returns a `200` status code with the following message:

```json
{ "message": "Deleted successfully" }
```

In any other case, it will return a corresponding error code and message.

### Possible `POST` error messages are:

| Error Message                                     | Explanation                                                                                                          |
| ------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| required fields are empty                         | Returned when one or more mandatory fields in the JSON payload are empty.                                            |
| invalid country code                              | Returned if the country code is not in ISO 3166-1 alpha-2 format.                                                    |
| invalid SWIFT code                                | Returned if the SWIFT code is not an 8 or 11-character alphanumeric string.                                          |
| invalid country name                              | Returned if the country name contains characters other than letters and spaces.                                      |
| country code does not match SWIFT country code    | Returned if the country code extracted from the provided SWIFT code does not match the one in the JSON payload.      |
| wrong country name                                | Returned if the provided country name does not match the one associated with the given country code in the database. |
| wrong country code                                | Returned if the provided country code does not match the one associated with the given country name in the database. |
| SWIFT code already exists                         | Returned if the SWIFT code is already present in the database.                                                       |
| headquarter flag does not match SWIFT code suffix | Returned if the headquarter flag does not match the value suggested by the last 3 characters of the SWIFT code.      |
| malformed JSON payload                            | Returned if the body of request is not in correct JSON format.                                                       |

### Possible `DELETE` error messages are:

| Error Message        | Explanation                                                                    |
| -------------------- | ------------------------------------------------------------------------------ |
| SWIFT code not found | Returned if the SWIFT code provided for deletion is not found in the database. |

### 'Wrong' vs 'invalid' in error codes

- Wrong: The data is semantically incorrect, even though it is in the correct format.
- Invalid: The data does not conform to the required format.
