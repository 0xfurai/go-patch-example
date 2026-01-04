# Golang HTTP Server with SQLite and GORM

A simple REST API server built with Go, featuring GET, POST, PUT, and PATCH endpoints connected to a SQLite database using GORM.

## Features

- HTTP server with RESTful endpoints
- SQLite database with GORM ORM
- CRUD operations (Create, Read, Update, Partial Update)
- Auto-migration of database schema
- Request validation using go-playground/validator
- DTOs (Data Transfer Objects) for all endpoints

## Prerequisites

- Go 1.16 or higher
- SQLite3 (usually comes pre-installed on most systems)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd golang-http-patch
```

2. Install dependencies:
```bash
go mod download
```

## Running the Server

Start the server:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### GET /users
Get all users.

**Response:**
```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
]
```

### GET /users/{id}
Get a single user by ID.

**Response:**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com"
}
```

### POST /users
Create a new user.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Validation Rules:**
- `name`: Required, minimum 2 characters, maximum 100 characters
- `email`: Required, must be a valid email address

**Response:**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Error Response (Validation Failed):**
```json
{
  "error": "Validation failed",
  "errors": [
    {
      "field": "email",
      "tag": "email",
      "message": "email must be a valid email address"
    }
  ]
}
```

### PUT /users/{id}
Update a user (full update - replaces all fields).

**Request Body:**
```json
{
  "name": "Jane Doe",
  "email": "jane@example.com"
}
```

**Validation Rules:**
- `name`: Required, minimum 2 characters, maximum 100 characters
- `email`: Required, must be a valid email address

**Response:**
```json
{
  "id": 1,
  "name": "Jane Doe",
  "email": "jane@example.com"
}
```

**Error Response (User Not Found):**
```
User not found
```

### PATCH /users/{id}
Partially update a user (only updates provided fields).

**Request Body:**
```json
{
  "name": "Jane Doe"
}
```

**Validation Rules:**
- `name`: Optional, if provided: minimum 2 characters, maximum 100 characters
- `email`: Optional, if provided: must be a valid email address
- At least one field must be provided

**Response:**
```json
{
  "id": 1,
  "name": "Jane Doe",
  "email": "john@example.com"
}
```

**Error Response (No Fields Provided):**
```
No fields to update
```

## Example Usage

### Create a user:
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

### Get all users:
```bash
curl http://localhost:8080/users
```

### Get a specific user:
```bash
curl http://localhost:8080/users/1
```

### Update a user (PUT):
```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com"}'
```

### Partially update a user (PATCH):
```bash
curl -X PATCH http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe"}'
```

## Database

The SQLite database file (`test.db`) will be created automatically in the project root directory when you first run the server. The database schema is automatically migrated using GORM's AutoMigrate feature.

## Project Structure

```
.
├── main.go          # Main application file with server and handlers
├── go.mod           # Go module dependencies
├── go.sum           # Go module checksums
├── test.db          # SQLite database (created on first run)
└── README.md        # This file
```

## Validation

All POST, PUT, and PATCH endpoints validate incoming requests using the `go-playground/validator` library. Validation errors are returned as JSON with detailed field-level error messages.

### Validation Rules

- **Name**: Required (POST/PUT), Optional (PATCH), 2-100 characters
- **Email**: Required (POST/PUT), Optional (PATCH), must be a valid email format

### Example Validation Error Response

```json
{
  "error": "Validation failed",
  "errors": [
    {
      "field": "name",
      "tag": "min",
      "message": "name must be at least 2 characters"
    },
    {
      "field": "email",
      "tag": "email",
      "message": "email must be a valid email address"
    }
  ]
}
```

## Dependencies

- [gorilla/mux](https://github.com/gorilla/mux) - HTTP router and URL matcher
- [gorm.io/gorm](https://gorm.io/) - The fantastic ORM library for Go
- [gorm.io/driver/sqlite](https://gorm.io/drivers/sqlite) - SQLite driver for GORM
- [go-playground/validator](https://github.com/go-playground/validator) - Struct validation library

