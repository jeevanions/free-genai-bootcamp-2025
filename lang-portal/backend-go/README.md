# Italian Language Learning Portal Backend

This is the backend service for the Italian Language Learning Portal, built with Go and SQLite.

## Prerequisites

- Go 1.21 or higher
- SQLite3

## Project Structure

The project follows clean architecture principles and is organized as follows:

```
.
├── cmd/                    # Application entry points
│   └── server/            # Main server binary
├── internal/              # Private application code
│   ├── api/              # API layer
│   ├── config/           # Configuration management
│   ├── domain/           # Business/domain logic
│   └── db/               # Database access
├── pkg/                  # Public shared packages
└── test/                 # Integration/E2E tests
```

## Getting Started

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up the database:
   ```bash
   goose -dir internal/db/migrations sqlite3 words.db up
   ```
4. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

The server will start on port 8080 by default. You can configure the port using the `PORT` environment variable.

## API Documentation

Once the server is running, you can access the Swagger documentation at:
```
http://localhost:8080/swagger/index.html
```

## Environment Variables

- `PORT`: Server port (default: 8080)
- `DB_PATH`: SQLite database path (default: words.db)
- `ENV_MODE`: Environment mode (development/production, default: development)

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Development

The project uses several tools to maintain code quality:

- `sqlc` for type-safe SQL
- `swaggo/swag` for API documentation
- `goose` for database migrations
- `zerolog` for structured logging

## License

MIT
