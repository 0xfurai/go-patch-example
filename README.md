# Golang HTTP Server with SQLite and GORM

A REST API server built with Go, featuring GET, POST, PUT, and PATCH endpoints connected to a SQLite database using GORM. This project demonstrates advanced PATCH operations with tri-state optional fields (unset, null, value) for partial updates.

## Features

- HTTP server with RESTful endpoints
- SQLite database with GORM ORM
- CRUD operations (Create, Read, Update, Partial Update)
- Advanced PATCH implementation with `patch.Optional[T]` type supporting three states:
  - **Unset**: Field not provided (ignored in update)
  - **Null**: Field explicitly set to null (removes/sets to null)
  - **Value**: Field provided with a value (updates to that value)
- Auto-migration of database schema
- Request validation using go-playground/validator with custom patch validators
- DTOs (Data Transfer Objects) for all endpoints
- Structured project layout with separate packages

## Prerequisites

- Go 1.24.3 or higher
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
    "email": "john@example.com",
    "age": 30,
    "phone": "+1234567890",
    "active": true,
    "bio": "Software developer",
    "role": "user",
    "score": 85.5
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
  "email": "john@example.com",
  "age": 30,
  "phone": "+1234567890",
  "active": true,
  "bio": "Software developer",
  "role": "user",
  "score": 85.5
}
```

**Error Response (User Not Found):**
```
User not found
```

### POST /users
Create a new user.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "age": 30,
  "phone": "+1234567890",
  "active": true,
  "bio": "Software developer",
  "role": "user",
  "score": 85.5
}
```

**Validation Rules:**
- `name`: Required, minimum 2 characters, maximum 100 characters
- `email`: Required, must be a valid email address (unique)
- `age`: Required, must be between 0 and 150
- `phone`: Optional, if provided: minimum 10 characters, maximum 20 characters
- `active`: Optional, defaults to `true`
- `bio`: Optional, maximum 500 characters
- `role`: Optional, must be one of: `admin`, `user`, `guest` (defaults to `user`)
- `score`: Optional, must be between 0 and 100 (defaults to 0)

**Response:**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "age": 30,
  "phone": "+1234567890",
  "active": true,
  "bio": "Software developer",
  "role": "user",
  "score": 85.5
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
Update a user (full update - replaces all fields). Note: Email is immutable and cannot be updated after creation.

**Request Body:**
```json
{
  "name": "Jane Doe",
  "age": 25,
  "phone": "+9876543210",
  "active": false,
  "bio": "Updated bio",
  "role": "admin",
  "score": 95.0
}
```

**Validation Rules:**
- `name`: Required, minimum 2 characters, maximum 100 characters
- `age`: Required, must be between 0 and 150
- `phone`: Optional, if provided: minimum 10 characters, maximum 20 characters
- `active`: Required
- `bio`: Optional, maximum 500 characters
- `role`: Required, must be one of: `admin`, `user`, `guest`
- `score`: Required, must be between 0 and 100
- `email`: Immutable (cannot be updated)

**Response:**
```json
{
  "id": 1,
  "name": "Jane Doe",
  "email": "john@example.com",
  "age": 25,
  "phone": "+9876543210",
  "active": false,
  "bio": "Updated bio",
  "role": "admin",
  "score": 95.0
}
```

**Error Response (User Not Found):**
```
User not found
```

### PATCH /users/{id}
Partially update a user (only updates provided fields). Supports three states for each field:
- **Unset**: Field not included in request (ignored, not updated)
- **Null**: Field explicitly set to `null` (removes/sets to null)
- **Value**: Field provided with a value (updates to that value)

**Request Body Examples:**

Update only name:
```json
{
  "name": "Jane Doe"
}
```

Update multiple fields:
```json
{
  "name": "Jane Doe",
  "age": 25,
  "active": false
}
```

Remove phone (set to null):
```json
{
  "phone": null
}
```

Mixed update (set some, remove others):
```json
{
  "name": "Jane Doe",
  "phone": null,
  "bio": null,
  "score": 90.0
}
```

**Validation Rules:**
- `name`: Optional, if provided: minimum 2 characters, maximum 100 characters
- `age`: Optional, if provided: must be between 0 and 150
- `phone`: Optional, if provided: minimum 10 characters, maximum 20 characters (can be set to null)
- `active`: Optional, can be set to null
- `bio`: Optional, if provided: maximum 500 characters (can be set to null)
- `role`: Optional, if provided: must be one of: `admin`, `user`, `guest` (can be set to null)
- `score`: Optional, if provided: must be between 0 and 100 (can be set to null)
- `email`: Immutable (cannot be updated)
- At least one field must be provided

**Response:**
```json
{
  "id": 1,
  "name": "Jane Doe",
  "email": "john@example.com",
  "age": 25,
  "phone": null,
  "active": false,
  "bio": null,
  "role": "user",
  "score": 90.0
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
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30,
    "phone": "+1234567890",
    "active": true,
    "bio": "Software developer",
    "role": "user",
    "score": 85.5
  }'
```

### Get all users:
```bash
curl http://localhost:8080/users
```

### Get a specific user:
```bash
curl http://localhost:8080/users/1
```

### Update a user (PUT - full update):
```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "age": 25,
    "phone": "+9876543210",
    "active": false,
    "bio": "Updated bio",
    "role": "admin",
    "score": 95.0
  }'
```

### Partially update a user (PATCH - update specific fields):
```bash
curl -X PATCH http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Doe", "age": 25}'
```

### Remove a field (PATCH - set to null):
```bash
curl -X PATCH http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"phone": null, "bio": null}'
```

## Database

The SQLite database file (`test.db`) will be created automatically in the project root directory when you first run the server. The database schema is automatically migrated using GORM's AutoMigrate feature.

## Project Structure

```
.
├── main.go              # Main application file with server setup and routing
├── go.mod               # Go module dependencies
├── go.sum               # Go module checksums
├── integration_test.go  # Integration tests for all endpoints
├── test.db              # SQLite database (created on first run)
├── README.md            # This file
├── database/
│   └── db.go            # Database initialization and connection
├── handlers/
│   ├── create_user.go   # POST /users handler
│   ├── get_user.go      # GET /users/{id} handler
│   ├── get_users.go     # GET /users handler
│   ├── patch_user.go    # PATCH /users/{id} handler
│   └── update_user.go   # PUT /users/{id} handler
├── models/
│   └── user.go          # User model and DTOs (CreateUserDTO, UpdateUserDTO, PatchUserDTO)
├── patch/
│   └── optional.go      # patch.Optional[T] type for tri-state PATCH operations
└── validation/
    ├── patchval.go      # Custom validators for patch.Optional types
    └── validator.go     # Validation setup and helper functions
```

## Validation

All POST, PUT, and PATCH endpoints validate incoming requests using the `go-playground/validator` library. Validation errors are returned as JSON with detailed field-level error messages.

### Validation Rules

#### POST /users (CreateUserDTO)
- **name**: Required, 2-100 characters
- **email**: Required, must be a valid email format (unique)
- **age**: Required, must be between 0 and 150
- **phone**: Optional, if provided: 10-20 characters
- **active**: Optional, defaults to `true`
- **bio**: Optional, maximum 500 characters
- **role**: Optional, must be one of: `admin`, `user`, `guest` (defaults to `user`)
- **score**: Optional, must be between 0 and 100 (defaults to 0)

#### PUT /users/{id} (UpdateUserDTO)
- **name**: Required, 2-100 characters
- **age**: Required, must be between 0 and 150
- **phone**: Optional, if provided: 10-20 characters
- **active**: Required
- **bio**: Optional, maximum 500 characters
- **role**: Required, must be one of: `admin`, `user`, `guest`
- **score**: Required, must be between 0 and 100
- **email**: Immutable (cannot be updated)

#### PATCH /users/{id} (PatchUserDTO)
Uses custom `opt` validator tag for `patch.Optional[T]` fields:
- **name**: Optional, if provided: 2-100 characters (can be set to null)
- **age**: Optional, if provided: must be between 0 and 150 (can be set to null)
- **phone**: Optional, if provided: 10-20 characters (can be set to null)
- **active**: Optional (can be set to null)
- **bio**: Optional, if provided: maximum 500 characters (can be set to null)
- **role**: Optional, if provided: must be one of: `admin`, `user`, `guest` (can be set to null)
- **score**: Optional, if provided: must be between 0 and 100 (can be set to null)
- **email**: Immutable (cannot be updated)

### Custom Patch Validators

The project includes custom validators for PATCH operations:
- **`opt`**: Validates optional fields only when they have a value (not unset or null)
  - Syntax: `validate:"opt=min=2;max=100"` (rules separated by semicolons)
  - Unset fields: validation passes
  - Null fields: validation passes
  - Value fields: validates with inner rules
- **`nonull`**: Prevents null values (fields can be unset or have a value, but not null)

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
    },
    {
      "field": "age",
      "tag": "gte",
      "message": "age must be 0 or greater"
    }
  ]
}
```

## Advanced PATCH Implementation

This project implements a sophisticated PATCH mechanism using the `patch.Optional[T]` generic type. This allows for true partial updates with three distinct states:

1. **Unset**: Field is not included in the JSON request → field is ignored (not updated)
2. **Null**: Field is explicitly set to `null` → field is set to null/removed
3. **Value**: Field is provided with a value → field is updated to that value

This is particularly useful for nullable fields like `phone` and `bio`, where you might want to:
- Leave the field unchanged (omit it)
- Remove the field (set to null)
- Update the field (provide a new value)

The implementation uses custom JSON unmarshaling and validation to handle these three states correctly.

## Testing

Run integration tests:
```bash
go test -v
```

The integration tests cover all endpoints and validation scenarios, including the advanced PATCH operations with unset, null, and value states.

## Dependencies

- [gorilla/mux](https://github.com/gorilla/mux) - HTTP router and URL matcher
- [gorm.io/gorm](https://gorm.io/) - The fantastic ORM library for Go
- [gorm.io/driver/sqlite](https://gorm.io/drivers/sqlite) - SQLite driver for GORM
- [go-playground/validator/v10](https://github.com/go-playground/validator) - Struct validation library with custom validators

