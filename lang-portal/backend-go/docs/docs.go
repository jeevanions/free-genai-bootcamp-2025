// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/dashboard/last_study_session": {
            "get": {
                "description": "Returns information about the most recent study session",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dashboard"
                ],
                "summary": "Get last study session",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.DashboardLastStudySession"
                        }
                    }
                }
            }
        },
        "/api/dashboard/quick-stats": {
            "get": {
                "description": "Returns quick overview statistics",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dashboard"
                ],
                "summary": "Get quick stats",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.DashboardQuickStats"
                        }
                    }
                }
            }
        },
        "/api/dashboard/study_progress": {
            "get": {
                "description": "Returns study progress statistics",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dashboard"
                ],
                "summary": "Get study progress",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.DashboardStudyProgress"
                        }
                    }
                }
            }
        },
        "/api/groups": {
            "get": {
                "description": "Returns a paginated list of groups",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "groups"
                ],
                "summary": "Get all groups",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 100,
                        "description": "Number of items per page",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "Offset for pagination",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.GroupListResponse"
                        }
                    }
                }
            }
        },
        "/api/groups/{id}": {
            "get": {
                "description": "Returns details about a specific group",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "groups"
                ],
                "summary": "Get group by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Group ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.GroupDetailResponse"
                        }
                    }
                }
            }
        },
        "/api/groups/{id}/study_sessions": {
            "get": {
                "description": "Returns a paginated list of study sessions for a specific group",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "groups"
                ],
                "summary": "Get study sessions for a group",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Group ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 100,
                        "description": "Number of items per page",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "Offset for pagination",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.GroupStudySessionsResponse"
                        }
                    }
                }
            }
        },
        "/api/groups/{id}/words": {
            "get": {
                "description": "Returns a paginated list of words in a specific group",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "groups"
                ],
                "summary": "Get words in a group",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Group ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 100,
                        "description": "Number of items per page",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "Offset for pagination",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.GroupWordsResponse"
                        }
                    }
                }
            }
        },
        "/api/study_activities": {
            "get": {
                "description": "Returns a list of available study activities",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "study_activities"
                ],
                "summary": "Get all study activities",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 100,
                        "description": "Number of items per page",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "Offset for pagination",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.StudyActivityListResponse"
                        }
                    }
                }
            }
        },
        "/api/study_activities/{id}": {
            "get": {
                "description": "Returns details about a specific study activity",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "study_activities"
                ],
                "summary": "Get study activity details",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Study Activity ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.StudyActivityResponse"
                        }
                    }
                }
            }
        },
        "/api/study_activities/{id}/launch": {
            "post": {
                "description": "Returns a list of study sessions for a specific activity\nLaunches a new study activity session for a specific group",
                "consumes": [
                    "application/json",
                    "application/json"
                ],
                "produces": [
                    "application/json",
                    "application/json"
                ],
                "tags": [
                    "study_activities",
                    "study_activities"
                ],
                "summary": "Launch a new study activity session",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Study Activity ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Study Activity ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Launch request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.LaunchStudyActivityRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.LaunchStudyActivityResponse"
                        }
                    }
                }
            }
        },
        "/api/study_activities/{id}/study_sessions": {
            "get": {
                "description": "Returns a list of study sessions for a specific activity\nLaunches a new study activity session for a specific group",
                "consumes": [
                    "application/json",
                    "application/json"
                ],
                "produces": [
                    "application/json",
                    "application/json"
                ],
                "tags": [
                    "study_activities",
                    "study_activities"
                ],
                "summary": "Launch a new study activity session",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Study Activity ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Study Activity ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Launch request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.LaunchStudyActivityRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.LaunchStudyActivityResponse"
                        }
                    }
                }
            }
        },
        "/api/words": {
            "get": {
                "description": "Returns a paginated list of words",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "words"
                ],
                "summary": "Get all words",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 100,
                        "description": "Number of items per page",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "Offset for pagination",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.WordListResponse"
                        }
                    }
                }
            }
        },
        "/api/words/{id}": {
            "get": {
                "description": "Returns details about a specific word",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "words"
                ],
                "summary": "Get word by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Word ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.WordResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.DashboardLastStudySession": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "group_id": {
                    "type": "integer"
                },
                "group_name": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "study_activity_id": {
                    "type": "integer"
                }
            }
        },
        "models.DashboardQuickStats": {
            "type": "object",
            "properties": {
                "study_streak_days": {
                    "type": "integer"
                },
                "success_rate": {
                    "type": "number"
                },
                "total_active_groups": {
                    "type": "integer"
                },
                "total_study_sessions": {
                    "type": "integer"
                }
            }
        },
        "models.DashboardStudyProgress": {
            "type": "object",
            "properties": {
                "total_available_words": {
                    "type": "integer"
                },
                "total_words_studied": {
                    "type": "integer"
                }
            }
        },
        "models.GroupDetailResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "stats": {
                    "$ref": "#/definitions/models.GroupStats"
                }
            }
        },
        "models.GroupListResponse": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.GroupResponse"
                    }
                },
                "pagination": {
                    "$ref": "#/definitions/models.PaginationResponse"
                }
            }
        },
        "models.GroupResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "word_count": {
                    "type": "integer"
                }
            }
        },
        "models.GroupStats": {
            "type": "object",
            "properties": {
                "total_word_count": {
                    "type": "integer"
                }
            }
        },
        "models.GroupStudySessionsResponse": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.StudySessionResponse"
                    }
                },
                "pagination": {
                    "$ref": "#/definitions/models.PaginationResponse"
                }
            }
        },
        "models.GroupWordsResponse": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.WordResponse"
                    }
                },
                "pagination": {
                    "$ref": "#/definitions/models.PaginationResponse"
                }
            }
        },
        "models.LaunchStudyActivityRequest": {
            "type": "object",
            "required": [
                "group_id"
            ],
            "properties": {
                "group_id": {
                    "type": "integer"
                }
            }
        },
        "models.LaunchStudyActivityResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "group_id": {
                    "type": "integer"
                },
                "study_activity_id": {
                    "type": "integer"
                },
                "study_session_id": {
                    "type": "integer"
                }
            }
        },
        "models.PaginationResponse": {
            "type": "object",
            "properties": {
                "current_page": {
                    "type": "integer"
                },
                "items_per_page": {
                    "type": "integer"
                },
                "total_items": {
                    "type": "integer"
                },
                "total_pages": {
                    "type": "integer"
                }
            }
        },
        "models.StudyActivityListResponse": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.StudyActivityResponse"
                    }
                },
                "pagination": {
                    "$ref": "#/definitions/models.PaginationResponse"
                }
            }
        },
        "models.StudyActivityResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "thumbnail_url": {
                    "type": "string"
                }
            }
        },
        "models.StudySessionResponse": {
            "type": "object",
            "properties": {
                "activity_name": {
                    "type": "string"
                },
                "correct_count": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "string"
                },
                "group_id": {
                    "type": "integer"
                },
                "group_name": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "words_count": {
                    "type": "integer"
                }
            }
        },
        "models.StudySessionsListResponse": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.StudySessionResponse"
                    }
                }
            }
        },
        "models.WordListResponse": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.WordResponse"
                    }
                },
                "pagination": {
                    "$ref": "#/definitions/models.PaginationResponse"
                }
            }
        },
        "models.WordResponse": {
            "type": "object",
            "properties": {
                "correct_count": {
                    "type": "integer"
                },
                "english": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "italian": {
                    "type": "string"
                },
                "parts": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "wrong_count": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "Italian Language Learning Portal API",
	Description:      "API for the Italian Language Learning Portal",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
