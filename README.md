# Awesome Project

A Go web application built with Gin framework, following clean architecture principles.

## Project Structure

```
.
├── cmd/                    # Application entry points
├── internal/              # Private application code
│   ├── delivery/         # API handlers and middleware
│   ├── domain/           # Business logic and models & Service layer
│      ├──models           #models
│      └──service         #service layer for endpoint | business logic
│   ├── service/          #mind future: later update, move service layer from domain to service folder         
│   ├── repository/       # Data access layer
│   ├── database/         # Database configuration
│   ├── usecase/          # Use case implementations
│   └── migrate/          # Migrate tables for current project
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
[Link share's Postman Collection](https://restless-flare-433229.postman.co/workspace/New-Team-Workspace~e0f4ed3a-5060-4895-bfb4-427ff650ec8b/collection/4175354-ae3c3b77-5c57-4cc1-9e53-ed22a5b8e7dc?action=share&creator=4175354)

## User
- `POST /login` - Login endpoint
  - Request body: `{ "email": "string", "password": "string" }`
  - Returns JWT token on success
- `POST /registration/user`
  - Request `{ "name": "Ryan Gosling","age": 32,"city": "Almaty","password": "dauren","email": "Ryan.gosling@example.com"}`
  - Returns response of created user `{"message": "User created successfully","user": { "created_at": "2025-05-20T14:46:35.607489+05:00","updated_at": "2025-05-20T14:46:35.607489+05:00","id": 2,"name": "Ryan Gosling","age": 32,"city": "Almaty","email": "Ryan.gosling@example.com"}}`
- `GET /protected/user/name/:id` - Get user information (protected route)
  - Requires JWT token in Authorization header
  - Returns user data for the specified ID

## News
- `GET /protected/news/all` - Get all news data (protected route)
  - Requires JWT token in Authorization header
  - Returns all news data
- `GET /protected/news/:id` - Get news data (protected route)
  - Requires JWT token in Authorization header
  - Returns news data for the specified ID

- `GET /` - Root endpoint
  - Returns a simple "Hello World" message

## Authentication

The application uses JWT (JSON Web Tokens) for authentication. Protected routes require a valid JWT token in the Authorization header:

```
Authorization: Bearer <your-token>
```
