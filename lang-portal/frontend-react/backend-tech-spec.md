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
  - id integer
  - italian string
  - english string
  - parts json
- words_groups - join table for words and groups many-to-many
  - id integer
  - word_id integer
  - group_id integer
- groups - thematic groups of words
  - id integer
  - name string
- study_sessions - records of study sessions grouping word_review_items
  - id integer
  - group_id integer
  - created_at datetime
  - study_activity_id integer
- study_activities - a specific study activity, linking a study session to group
  - id integer
  - study_session_id integer
  - group_id integer
  - created_at datetime
- word_review_items - a record of word practice, determining if the word was correct or not
  - word_id integer
  - study_session_id integer
  - correct boolean
  - created_at datetime

### Relationships

* word belongs to groups through  word_groups
* group belongs to words through word_groups
* session belongs to a group
* session belongs to a study_activity
* session has many word_review_items
* word_review_item belongs to a study_session
* word_review_item belongs to a word

### Design Notes

* All tables use auto-incrementing primary keys
* Timestamps are automatically set on creation where applicable
* Foreign key constraints maintain referential integrity
* JSON storage for word parts allows flexible component storage
* Counter cache on groups.words_count optimizes word counting queries


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
  "description": "Practice your vocabulary with flashcards"
}
```

### GET /api/dashboard/study-sessions

Returns a paginated list of all study sessions across all activities.

#### Query Parameters
- `page`: Page number (default: 1)
- `per_page`: Items per page (default: 10)

#### JSON Response
```json
{
  "items": [
    {
      "id": 123,
      "activity_name": "Vocabulary Quiz",
      "group_name": "Basic Greetings",
      "review_items_count": 10,
      "correct_count": 8,
      "wrong_count": 2,
      "start_time": "2025-02-08T17:20:23-05:00",
      "end_time": "2025-02-08T17:25:23-05:00",
      "duration_seconds": 300
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 45,
    "per_page": 10
  }
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
      "japanese": "こんにちは",
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
  "japanese": "こんにちは",
  "romaji": "konnichiwa",
  "english": "hello",
  "stats": {
    "correct_count": 5,
    "wrong_count": 2
  },
  "groups": [
    {
      "id": 1,
      "name": "Basic Greetings"
    }
  ]
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
      "japanese": "こんにちは",
      "romaji": "konnichiwa",
      "english": "hello",
      "correct_count": 5,
      "wrong_count": 2
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
      "japanese": "こんにちは",
      "romaji": "konnichiwa",
      "english": "hello",
      "correct_count": 5,
      "wrong_count": 2
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
#### JSON Response
```json
{
  "success": true,
  "message": "Study history has been reset"
}
```

### POST /api/full_reset
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

List of tasks required to build and test the backend application

### Initialize Database

This taks will initialize the sqllite database called `words.db` in the root of the project.

### Migrate databse

Task will run a series of migration sql files on the dtabase.

Migration live in the `migrations` folder. The migration files will be run in order of their file name. The file names should look like this

```
0001_init.sql
0002_create_words_table.sql
```

### Seed data

This will import a json files and transform them into target data of our database.

All Seed files live in the `seeds` folder.

### Reset

Resets the database by removing the database file

### ResetAndSeed

Resets the database and runs the seed data

### TestDB 

Initializes the test database with test data

### UnitTest

Runs unit tests

### Tidy

Runs go tidy command 

### Lint

Runs go lint command

### Docs

Updates swagger docs

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
