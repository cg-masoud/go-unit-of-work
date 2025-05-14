# Go Unit of Work Pattern Example

A simple Go project demonstrating the Unit of Work pattern for managing database transactions and ensuring data consistency.

## Project Structure

```
.
├── db/             # Database connection and Unit of Work implementation
├── handler/        # HTTP request handlers
├── middleware/     # HTTP middleware (Unit of Work middleware)
├── model/          # Domain models
├── repository/     # Data access layer
└── service/        # Business logic layer
```

## Key Features

- Unit of Work pattern implementation for transaction management
- Clean architecture with separation of concerns
- Repository pattern for data access
- Middleware for automatic transaction handling
- RESTful API endpoints for order management

## Getting Started

1. Clone the repository
2. Install dependencies: `go mod download`
3. Run the application: `go run main.go`

## API Endpoints

- `DELETE /orders/:id` - Delete order and order items