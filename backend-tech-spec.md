# Technical Specification: Italian Language Learning Portal Backend

## 1. Overview

### 1.1 Business Objective

The objective is to build a platform that will allow students to learn Italian by practicing vocabulary, grammar, and pronunciation.

The platform will have the following core capabilities:
1. **Vocabulary Management**: Central repository for 5000+ Italian words with linguistic metadata
2. **Learning Record Store**: Record practice performance with spaced repetition statistics
3. **Activity Launchpad**: Unified interface for launching learning exercises

## 2. System Architecture

### 2.1 Tech Stack

| Component       | Technology       | Rationale                                                                 |
|-----------------|------------------|---------------------------------------------------------------------------|
| Language        | Go 1.21+         | Strong concurrency support, excellent performance characteristics        |
| Web Framework   | Gin v1.9.1       | High-performance HTTP framework with middleware ecosystem                |
| Database        | SQLite3          | Single-file storage, ACID compliance, suitable for single-user prototype |
| Task Runner     | Mage v1.15.0     | Go-native build tool with declarative task definitions                   |
| ORM             | sqlc v1.24.0     | Type-safe SQL to Go code generation                                      |
| Migrations      | goose v3.19.0    | Robust database migration management                                     |

API Format: RESTful JSON endpoints

Auth: None (single-user system)

### 2.2 Architectural Diagram

#### 2.2.1 Component Diagram

[Frontend] 
  ↔ [Gin HTTP Server]
    ↔ [SQLite Database]
    ↔ [Business Logic]
    ↔ [Session Manager]

### 2.3 Library Dependencies

| Package             | Version | Purpose                              |
|---------------------|---------|--------------------------------------|
| github.com/gin-gonic/gin | v1.9.1  | HTTP routing and middleware          |
| modernc.org/sqlite  | v1.29.0 | Pure-Go SQLite implementation        |
| github.com/stretchr/testify | v1.9.0 | Testing framework                    |
| github.com/rs/zerolog | v1.32.0 | Structured logging                   |
| github.com/prometheus/client_golang | v1.19.0 | Metrics collection                  |

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

## 3. Database Definition

### 3.1 Schema Definition

```sql
-- Core Vocabulary with Italian-specific features
CREATE TABLE words (
    id INTEGER PRIMARY KEY,
    italian TEXT NOT NULL,                -- Base form of the word
    english TEXT NOT NULL,                -- English translation
    parts_of_speech TEXT NOT NULL,        -- noun, verb, adjective, etc.
    gender TEXT,                          -- masculine/feminine for nouns
    number TEXT,                          -- singular/plural forms
    difficulty_level INTEGER NOT NULL,    -- 1-5 scale for progression
    verb_conjugation JSON,                -- Store all tenses for verbs
    notes TEXT,                           -- Usage notes, idioms
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Words to Groups Many-to-Many Relationship
CREATE TABLE words_groups (
    id INTEGER PRIMARY KEY,
    word_id INTEGER REFERENCES words(id),
    group_id INTEGER REFERENCES groups(id),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Thematic Groups (e.g., Food, Travel, Business)
CREATE TABLE groups (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    difficulty_level INTEGER NOT NULL,    -- Match with word difficulty
    category TEXT NOT NULL,               -- grammar/thematic/situational
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Enhanced Study Sessions with Metrics
CREATE TABLE study_sessions (
    id INTEGER PRIMARY KEY,
    group_id INTEGER REFERENCES groups(id),
    study_activity_id INTEGER REFERENCES study_activities(id),
    total_words INTEGER NOT NULL,         -- Words attempted
    correct_words INTEGER NOT NULL,       -- Successful attempts
    duration_seconds INTEGER NOT NULL,    -- Session length
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Activity Types for Italian Learning
CREATE TABLE study_activities (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,                   -- vocabulary/grammar/pronunciation
    requires_audio BOOLEAN NOT NULL,      -- For pronunciation exercises
    difficulty_level INTEGER NOT NULL,    -- Progressive difficulty
    instructions TEXT NOT NULL,           -- Activity guidelines
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Word Review Records
CREATE TABLE word_review_items (
    id INTEGER PRIMARY KEY,
    word_id INTEGER REFERENCES words(id),
    study_session_id INTEGER REFERENCES study_sessions(id),
    correct BOOLEAN NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 3.2 Index Strategy

```sql
-- Words table indexes
CREATE INDEX idx_words_italian ON words(italian);         -- Fast word lookup
CREATE INDEX idx_words_difficulty ON words(difficulty_level);  -- For progression filtering
CREATE INDEX idx_words_parts_speech ON words(parts_of_speech);  -- Filter by word type

-- Words Groups indexes
CREATE INDEX idx_words_groups_word ON words_groups(word_id);    -- Word relationship lookups
CREATE INDEX idx_words_groups_group ON words_groups(group_id);   -- Group relationship lookups

-- Groups indexes
CREATE INDEX idx_groups_category ON groups(category);      -- Group filtering
CREATE INDEX idx_groups_difficulty ON groups(difficulty_level);  -- Difficulty filtering

-- Study Sessions indexes
CREATE INDEX idx_study_sessions_group ON study_sessions(group_id);  -- Group performance tracking
CREATE INDEX idx_study_sessions_created ON study_sessions(created_at);  -- Time-based queries
CREATE INDEX idx_study_sessions_metrics ON study_sessions(correct_words, total_words);  -- Performance analysis

-- Study Activities indexes
CREATE INDEX idx_study_activities_type ON study_activities(type);  -- Activity type filtering
CREATE INDEX idx_study_activities_difficulty ON study_activities(difficulty_level);  -- Difficulty filtering

-- Word Review Items indexes
CREATE INDEX idx_word_review_performance ON word_review_items(word_id, correct);  -- Word success tracking
CREATE INDEX idx_word_review_session ON word_review_items(study_session_id);  -- Session lookups
CREATE INDEX idx_word_review_created ON word_review_items(created_at);  -- Time-based analysis
```

## 4. API Endpoints

### 4.1 Dashboard Endpoints

**GET /api/dashboard/last_study_session**
```json
{
  "id": 123,
  "group_id": 456,
  "created_at": "2024-03-10T15:30:00Z",
  "study_activity_id": 789,
  "group_name": "Basic Verbs"
}
```

**GET /api/dashboard/study_progress**
```json
{
  "total_words_studied": 150,
  "total_available_words": 500,
  "groups_progress": [
    {
      "group_id": 1,
      "group_name": "Basic Verbs",
      "words_mastered": 45,
      "total_words": 100
    }
  ]
}
```

**GET /api/dashboard/quick-stats**
```json
{
  "success_rate": 80.0,
  "total_study_sessions": 4,
  "total_active_groups": 3,
  "study_streak_days": 4
}
```

### 4.2 Words Management

**GET /api/words**
```json
{
  "items": [
    {
      "id": 1,
      "italian": "ciao",
      "english": "hello",
      "parts": {
        "type": "greeting",
        "usage": "informal"
      },
      "parts_of_speech": "interjection",
      "difficulty_level": 1,
      "notes": "Common informal greeting",
      "stats": {
        "correct_count": 5,
        "incorrect_count": 2,
        "last_reviewed": "2024-03-10T15:30:00Z"
      }
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

**GET /api/words/:id**
```json
{
  "id": 1,
  "italian": "ciao",
  "english": "hello",
  "parts": {
    "type": "greeting",
    "usage": "informal"
  },
  "parts_of_speech": "interjection",
  "difficulty_level": 1,
  "notes": "Common informal greeting"
}
```

### 4.3 Groups Management

**GET /api/groups**
```json
{
  "items": [
    {
      "id": 1,
      "name": "Basic Greetings",
      "description": "Essential Italian greetings",
      "difficulty_level": 1,
      "category": "thematic",
      "word_count": 20,
      "completion_rate": 75.5
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

**GET /api/groups/:id/words**
```json
{
  "items": [
    {
      "id": 1,
      "italian": "ciao",
      "english": "hello",
      "stats": {
        "correct_count": 5,
        "incorrect_count": 2,
        "last_reviewed": "2024-03-10T15:30:00Z"
      }
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 20,
    "items_per_page": 20
  }
}
```

### 4.4 Study Sessions

**POST /api/study_sessions**
```json
Request:
{
  "group_id": 123,
  "study_activity_id": 456
}

Response:
{
  "id": 789,
  "group_id": 123,
  "study_activity_id": 456,
  "created_at": "2024-03-10T15:30:00Z"
}
```

**GET /api/study_sessions/:id**
```json
{
  "id": 789,
  "group_id": 123,
  "study_activity_id": 456,
  "total_words": 20,
  "correct_words": 15,
  "duration_seconds": 300,
  "created_at": "2024-03-10T15:30:00Z"
}
```

### 4.5 Word Reviews

**POST /api/study_sessions/:session_id/reviews**
```json
Request:
{
  "word_id": 1,
  "correct": true
}

Response:
{
  "id": 123,
  "word_id": 1,
  "study_session_id": 789,
  "correct": true,
  "created_at": "2024-03-10T15:30:00Z"
}
```

**GET /api/study_sessions/:session_id/reviews**
```json
{
  "items": [
    {
      "id": 123,
      "word_id": 1,
      "correct": true,
      "created_at": "2024-03-10T15:30:00Z",
      "word": {
        "italian": "ciao",
        "english": "hello"
      }
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 20,
    "items_per_page": 20
  }
}
```

### 4.6 Study Activities

**GET /api/study_activities**
```json
{
  "items": [
    {
      "id": 1,
      "type": "vocabulary",
      "requires_audio": false,
      "difficulty_level": 1,
      "instructions": "Match the Italian words with their English translations"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 5,
    "items_per_page": 20
  }
}
```

### 4.7 System Management

**POST /api/reset**
```json
Request:
{
  "reset_type": "study_history"  // or "full_system"
}

Response:
{
  "success": true,
  "message": "Study history has been reset",
  "reset_type": "study_history",
  "timestamp": "2024-03-10T15:30:00Z"
}
```

## 5. Implementation Details

### 5.1 Database Connection
```go
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