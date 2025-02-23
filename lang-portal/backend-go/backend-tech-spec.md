# Technical Specification: Italian Language Learning Portal Backend

## 1. Overview

### 1.1 Business Objective

A language learning school wants to build a prototype of learning portal which will act as three things:
* Inventory of possible vocabulary that can be learned
* Act as a  Learning record store (LRS), providing correct and wrong score on practice vocabulary
* A unified launchpad to launch different learning apps


## 2. System Architecture

### 2.1 Technical Stack

| Component       | Technology       | Rationale                                                                 |
|-----------------|------------------|---------------------------------------------------------------------------|
| Language        | Go 1.21+         | Strong concurrency support, excellent performance characteristics        |
| Web Framework   | Gin v1.9.1       | High-performance HTTP framework with middleware ecosystem                |
| Database        | SQLite3          | Single-file storage, ACID compliance, suitable for single-user prototype |
| Task Runner     | Mage v1.15.0     | Go-native build tool with declarative task definitions                   |
| ORM             | sqlc v1.24.0     | Type-safe SQL to Go code generation                                      |
| Migrations      | goose v3.19.0    | Robust database migration management                                     |
| API Docs        | swaggo/swag      | OpenAPI/Swagger documentation generation                                 |

API Format: RESTful JSON endpoints

Auth: None (single-user system)

Enable swagger UI so we can test the API from the browser.

### 2.3 Library Dependencies

| Package             | Version | Purpose                              |
|---------------------|---------|--------------------------------------|
| github.com/gin-gonic/gin | v1.9.1  | HTTP routing and middleware          |
| modernc.org/sqlite  | v1.29.0 | Pure-Go SQLite implementation        |
| github.com/stretchr/testify | v1.9.0 | Testing framework                    |
| github.com/rs/zerolog | v1.32.0 | Structured logging                   |
| github.com/swaggo/swag | v1.16.3 | Swagger documentation generation     |
| github.com/swaggo/gin-swagger | v1.6.0 | Swagger UI for Gin                   |

### 2.4 Backend project structure

```
.
├── cmd/                    # Application entry points
│   └── server/            # Main server binary
│       └── main.go        # Server initialization
├── internal/              # Private application code
│   ├── api/              # API layer
│   │   ├── handlers/     # HTTP request handlers
│   │   ├── middleware/   # HTTP middleware
│   │   └── router/       # Route definitions
│   ├── config/           # Configuration management
│   ├── domain/           # Business/domain logic
│   │   ├── models/       # Domain models
│   │   └── services/     # Business logic services
│   ├── db/               # Database access
│   │   ├── migrations/   # SQL migrations
│   │   ├── queries/      # SQL queries (for sqlc)
│   │   └── repository/   # Data access layer
│   └── pkg/              # Internal shared packages
│       ├── logger/       # Logging utilities
│       └── validator/    # Input validation
├── pkg/                  # Public shared packages
├── scripts/              # Build/deployment scripts
├── test/                 # Integration/E2E tests
├── .env.example         # Environment template
├── .gitignore
├── go.mod
├── go.sum
├── magefile.go          # Build automation
└── README.md
```

Key aspects of this structure:
- `cmd/`: Contains the main application entry points
- `internal/`: Private application code not meant for external use
- `pkg/`: Shared code that could be used by other projects
- Clear separation between API, domain logic, and data access layers
- Follows Go project layout conventions and clean architecture principles

## 3. Database schema

Our database will be a single sqlite database called `words.db` that will be in the root of the project folder of `backend_go`

We have the following tables:
- words - stored vocabulary words
  - id integer PRIMARY KEY AUTOINCREMENT
  - italian string NOT NULL
  - english string NOT NULL
  - parts json
  - created_at datetime DEFAULT CURRENT_TIMESTAMP

- words_groups - join table for words and groups many-to-many
  - id integer PRIMARY KEY AUTOINCREMENT
  - word_id integer NOT NULL REFERENCES words(id) ON DELETE CASCADE
  - group_id integer NOT NULL REFERENCES groups(id) ON DELETE CASCADE
  - created_at datetime DEFAULT CURRENT_TIMESTAMP
  - Indexes: word_id, group_id

- groups - thematic groups of words
  - id integer PRIMARY KEY AUTOINCREMENT
  - name string NOT NULL
  - words_count integer DEFAULT 0  # Counter cache for optimization
  - created_at datetime DEFAULT CURRENT_TIMESTAMP

- study_sessions - records of study sessions grouping word_review_items
  - id integer PRIMARY KEY AUTOINCREMENT
  - group_id integer NOT NULL REFERENCES groups(id) ON DELETE CASCADE
  - study_activity_id integer NOT NULL REFERENCES study_activities(id) ON DELETE CASCADE
  - created_at datetime DEFAULT CURRENT_TIMESTAMP
  - Indexes: group_id

Note: While the initial design considered separate start_time and end_time fields, the implementation uses created_at to track when the session started, as the current prototype focuses on simple session tracking without duration.

- study_activities - a specific study activity that can be launched with a group
  - id integer PRIMARY KEY AUTOINCREMENT
  - name string NOT NULL
  - thumbnail_url string
  - launch_url string
  - description string
  - created_at datetime DEFAULT CURRENT_TIMESTAMP

- word_review_items - a record of word practice, determining if the word was correct or not
  - id integer PRIMARY KEY AUTOINCREMENT
  - word_id integer NOT NULL REFERENCES words(id) ON DELETE CASCADE
  - study_session_id integer NOT NULL REFERENCES study_sessions(id) ON DELETE CASCADE
  - correct boolean NOT NULL
  - created_at datetime DEFAULT CURRENT_TIMESTAMP
  - Indexes: word_id, study_session_id

### Relationships

* word belongs to groups through  word_groups
* group belongs to words through word_groups
* session belongs to a group
* session belongs to a study_activity
* session has many word_review_items
* word_review_item belongs to a study_session
* word_review_item belongs to a word

### Design Notes

* All tables use auto-incrementing primary keys (INTEGER PRIMARY KEY AUTOINCREMENT)
* Timestamps are automatically set on creation using DEFAULT CURRENT_TIMESTAMP
* Foreign key constraints with ON DELETE CASCADE maintain referential integrity
* JSON storage for word parts allows flexible component storage
* Counter cache on groups.words_count optimizes word counting queries
* Performance optimized with indexes on frequently queried foreign keys
* NOT NULL constraints on required fields ensure data integrity


## 4. API Endpoints

### GET /api/dashboard/last_study_session
Returns information about the most recent study session.

#### JSON Response
```json
{
  "id": 123,
  "group_id": 456,
  "created_at": "2025-02-08T17:20:23-05:00",
  "study_activity_id": 789,
  "group_id": 456,
  "group_name": "Basic Greetings"
}
```

### GET /api/dashboard/study_progress
Returns study progress statistics.
Please note that the frontend will determine progress bar basedon total words studied and total available words.

#### JSON Response

```json
{
  "total_words_studied": 3,
  "total_available_words": 124,
}
```

### GET /api/dashboard/quick-stats

Returns quick overview statistics.

#### JSON Response
```json
{
  "success_rate": 80.0,
  "total_study_sessions": 4,
  "total_active_groups": 3,
  "study_streak_days": 4
}
```

### GET /api/study_activities/:id

#### JSON Response
```json
{
  "id": 1,
  "name": "Vocabulary Quiz",
  "thumbnail_url": "https://example.com/thumbnail.jpg",
  "description": "Practice your vocabulary with flashcards",
  "launch_url": "https://example.com/quiz/launch"
}
```

### POST /api/study_activities/:id/launch

Launches a new study activity session for a specific group.

#### Request Body
```json
{
  "group_id": 123
}
```

#### JSON Response
```json
{
  "study_session_id": 456,
  "study_activity_id": 789,
  "group_id": 123,
  "created_at": "2025-02-21T22:28:24Z"
}
```

### GET /api/study_activities/:id/study_sessions

- pagination with 100 items per page

```json
{
  "items": [
    {
      "id": 123,
      "activity_name": "Vocabulary Quiz",
      "group_name": "Basic Greetings",
      "start_time": "2025-02-08T17:20:23-05:00",
      "end_time": "2025-02-08T17:30:23-05:00",
      "review_items_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 100,
    "items_per_page": 20
  }
}
```

### POST /api/study_activities

#### Request Params
- group_id integer
- study_activity_id integer

#### JSON Response
{
  "id": 124,
  "group_id": 123
}

### GET /api/words

- pagination with 100 items per page

#### JSON Response
```json
{
  "items": [
    {
      "italian": "こんにちは",
      "romaji": "konnichiwa",
      "english": "hello",
      "correct_count": 5,
      "wrong_count": 2
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 500,
    "items_per_page": 100
  }
}
```

### GET /api/words/:id
#### JSON Response
```json
{
      "italian": "sorella",
      "english": "sister",
      "parts": {
        "type": "noun",
        "gender": "feminine",
        "plural": "sorelle"
      }
    }
```

### GET /api/groups
- pagination with 100 items per page
#### JSON Response
```json
{
  "items": [
    {
      "id": 1,
      "name": "Basic Greetings",
      "word_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 10,
    "items_per_page": 100
  }
}
```

### GET /api/groups/:id
#### JSON Response
```json
{
  "id": 1,
  "name": "Basic Greetings",
  "stats": {
    "total_word_count": 20
  }
}
```

### GET /api/groups/:id/words
#### JSON Response
```json
{
  "items": [
    {
      "id": 1,
      "italian": "ciao",
      "english": "hello",
      "parts": null,
      "correct_count": 0,
      "wrong_count": 0
    },
    {
      "id": 2,
      "italian": "buongiorno",
      "english": "good morning",
      "parts": null,
      "correct_count": 0,
      "wrong_count": 0
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 7,
    "items_per_page": 100
  }
}
```

### GET /api/groups/:id/study_sessions
#### JSON Response
```json
{
  "items": [
    {
      "id": 123,
      "activity_name": "Vocabulary Quiz",
      "group_name": "Basic Greetings",
      "start_time": "2025-02-08T17:20:23-05:00",
      "end_time": "2025-02-08T17:30:23-05:00",
      "review_items_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 5,
    "items_per_page": 100
  }
}
```

### GET /api/study_sessions
- pagination with 100 items per page
#### JSON Response
```json
{
  "items": [
    {
      "id": 123,
      "activity_name": "Vocabulary Quiz",
      "group_name": "Basic Greetings",
      "start_time": "2025-02-08T17:20:23-05:00",
      "end_time": "2025-02-08T17:30:23-05:00",
      "review_items_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 100,
    "items_per_page": 100
  }
}
```

### GET /api/study_sessions/:id
#### JSON Response
```json
{
  "id": 123,
  "activity_name": "Vocabulary Quiz",
  "group_name": "Basic Greetings",
  "start_time": "2025-02-08T17:20:23-05:00",
  "end_time": "2025-02-08T17:30:23-05:00",
  "review_items_count": 20
}
```

### GET /api/study_sessions/:id/words
- pagination with 100 items per page
#### JSON Response
```json
{
  "items": [
    {
      "id": 123,
      "activity_name": "Vocabulary Quiz",
      "group_name": "Basic Greetings",
      "start_time": "2025-02-08T17:20:23-05:00",
      "end_time": "2025-02-08T17:30:23-05:00",
      "review_items_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 20,
    "items_per_page": 100
  }
}
```

### POST /api/reset_history
Resets all study history data while preserving word and group data. This includes:
- Deleting all study sessions
- Deleting all word review items
- Resetting word correct/wrong counts

#### Description
Use this endpoint when you want to clear study progress but keep the vocabulary and group structure intact.

#### JSON Response
```json
{
  "success": true,
  "message": "Study history has been reset"
}
```

### POST /api/full_reset
Performs a complete system reset, including:
- Dropping all tables
- Recreating database schema
- Reseeding initial data

#### Description
Use this endpoint with caution as it will completely reset the system to its initial state. This is useful for:
- Development and testing
- Resetting a demo environment
- Troubleshooting data issues

#### JSON Response
```json
{
  "success": true,
  "message": "System has been fully reset"
}
```

### POST /api/study_sessions/:id/words/:word_id/review
#### Request Params
- id (study_session_id) integer
- word_id integer
- correct boolean

#### Request Payload
```json
{
  "correct": true
}
```

#### JSON Response
```json
{
  "success": true,
  "word_id": 1,
  "study_session_id": 123,
  "correct": true,
  "created_at": "2025-02-08T17:33:07-05:00"
}
```
## 5. Mage build tasks

List of tasks required to build and test the backend application. The implementation provides the following tasks:

### Database Tasks

#### db:migrate
Runs database migrations using the goose migration tool. Goose provides versioned migrations with up/down capabilities.

#### db:reset
Resets the database by dropping all tables and removing the database file.

#### db:seed
Seeds the database using the Go-based seeder package with JSON files from the `seeds` folder.

#### db:status
Checks and displays the current status of database migrations.

### Development Tasks

#### dev
Starts the application in development mode with hot reload.

#### run
Runs the application in normal mode.

#### install
Installs project dependencies.

#### all
Runs a complete build cycle including tests and linting.

#### build
Builds the application binary.

#### clean
Cleans build artifacts and temporary files.

### Testing and Quality Tasks

#### test
Runs the test suite.

#### coverage
Runs tests with coverage reporting.

#### lint
Runs the Go linter to ensure code quality.

### Code Generation Tasks

#### generate
Runs code generation tools.

#### swagger
Generates Swagger/OpenAPI documentation from code annotations.

### Task Execution Examples

```bash
# Database tasks
mage db:migrate    # Run migrations
mage db:status     # Check migration status
mage db:reset      # Reset database
mage db:seed       # Seed data

# Development tasks
mage dev          # Start development server
mage build        # Build application

# Testing tasks
mage test         # Run tests
mage coverage     # Run tests with coverage
mage lint         # Run linter

# Documentation
mage swagger      # Generate API docs
```

All tasks are designed to be composable and can be run individually or as part of a larger build process using the `all` task.

## 5. Implementation Details

### 5.1 Database Connection
```go
// internal/models/db.go
// internal/models/db.go
func NewDB() (*sql.DB, error) {
    return sql.Open("sqlite", "./words.db")
}
```

### 5.2 Request Validation
```go
// internal/handlers/sessions.go
type ReviewRequest struct {
    Correct bool `json:"correct" binding:"required"`
}

func SubmitReview(c *gin.Context) {
    var req ReviewRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid review data"})
        return
    }
}
```

## 6. Development Workflow

### 6.1 Mage Tasks
```go
// magefile.go
func Migrate() error {
    return sh.Run("sqlite3", "words.db", ".read db/migrations/0001_init.sql")
}

func Seed() error {
    return sh.Run("sqlite3", "words.db", ".read db/seeds/initial_data.sql")
}
```

### 6.2 Setup Guide
```bash
# Initialize system
go mod init
mage init
mage migrate
mage seed

# Start server
go run cmd/server/main.go
```

## 7. Error Handling

### 7.1 Standard Error Response
```json
{
  "error": "Invalid session ID",
  "code": "ERR-400-002",
  "documentation": "/api/docs#error-codes"
}
```

### 7.2 Common Status Codes
- 400: Invalid request format
- 404: Resource not found
- 500: Database connection error
- 422: Validation error

### 7.3 Error Logging
All errors will be logged using zerolog with contextual information including:
- Request ID
- User IP
- Endpoint
- Error stack trace 

## 8. Testing Strategy

### 8.1 Testing Framework Stack

| Component          | Technology                | Purpose                                          |
|-------------------|---------------------------|--------------------------------------------------|
| Testing Framework | testify v1.9.0           | Assertions and test suite organization           |
| HTTP Testing      | httptest (stdlib)        | HTTP handler testing                             |
| Mock Generation   | mockery v2.40.1          | Interface mocking for unit tests                 |
| Coverage Tool     | go test -cover           | Code coverage reporting                          |
| Performance Tests | k6 v0.49.0               | Load and performance testing                     |

### 8.2 Test Categories

1. **Unit Tests**
   - Location: Next to the source files (`*_test.go`)
   - Naming: `TestXxx` for test functions
   - Coverage target: 80% minimum
   - Mock external dependencies
   - Focus on single component behavior

### 8.3 Testing Standards

1. **Test Structure (AAA Pattern)**
   ```go
   func TestSomething(t *testing.T) {
       // Arrange
       expected := "expected result"
       sut := NewSystemUnderTest()

       // Act
       result, err := sut.DoSomething()

       // Assert
       assert.NoError(t, err)
       assert.Equal(t, expected, result)
   }
   ```

2. **Table-Driven Tests**
   ```go
   func TestOperation(t *testing.T) {
       tests := []struct {
           name     string
           input    string
           expected string
           wantErr  bool
       }{
           {
               name:     "valid input",
               input:    "test",
               expected: "result",
               wantErr:  false,
           },
           // More test cases...
       }

       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               // Test implementation
           })
       }
   }
   ```

3. **Mock Generation**
   ```go
   //go:generate mockery --name=Repository --output=mocks --outpkg=mocks --case=snake
   type Repository interface {
       FindByID(id string) (*Entity, error)
   }
   ```

### 8.4 Test Execution

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -run TestSpecificFunction

# Run performance tests
k6 run test/performance/scenario.js
```
