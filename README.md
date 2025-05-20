# Awesome Project

A Go web application built with Gin framework, following clean architecture principles.

## Project Structure

```
.
├── cmd/                    # Application entry points
├── internal/              # Private application code
│   ├── delivery/         # API handlers and middleware
│   ├── domain/           # Business logic and models
│   ├── service/          # Service layer
│   ├── repository/       # Data access layer
│   ├── database/         # Database configuration
│   └── usecase/          # Use case implementations
├── main.go               # Main application entry point
├── go.mod               # Go module file
└── go.sum               # Go module checksum file
```

## Dependencies

- Gin (Web Framework)
- GORM (ORM)
- PostgreSQL (Database)
- JWT (Authentication)

## Getting Started

1. Make sure you have Go 1.24 or later installed
2. Clone the repository
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Set up your PostgreSQL database
5. Run the application:
   ```bash
   go run main.go
   ```

The server will start on `localhost:8080`

## API Endpoints

- `POST /login` - Login endpoint
  - Request body: `{ "email": "string", "password": "string" }`
  - Returns JWT token on success
- `POST /registration/user`
  - Request `{ "name": "Ryan Gosling","age": 32,"city": "Almaty","password": "dauren","email": "Ryan.gosling@example.com"}`
  - Returns response of created user `{"message": "User created successfully","user": { "created_at": "2025-05-20T14:46:35.607489+05:00","updated_at": "2025-05-20T14:46:35.607489+05:00","id": 2,"name": "Ryan Gosling","age": 32,"city": "Almaty","email": "Ryan.gosling@example.com"}}`
- `GET /protected/user/name/:id` - Get user information (protected route)
  - Requires JWT token in Authorization header
  - Returns user data for the specified ID

- `GET /` - Root endpoint
  - Returns a simple "Hello World" message

## Authentication

The application uses JWT (JSON Web Tokens) for authentication. Protected routes require a valid JWT token in the Authorization header:

```
Authorization: Bearer <your-token>
```
