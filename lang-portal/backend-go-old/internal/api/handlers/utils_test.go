package handlers

import (
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestParseIDParam(t *testing.T) {
    gin.SetMode(gin.TestMode)

    tests := []struct {
        name        string
        paramValue  string
        expectedID  int64
        expectError bool
    }{
        {
            name:        "Valid ID",
            paramValue:  "123",
            expectedID:  123,
            expectError: false,
        },
        {
            name:        "Invalid ID - Not a number",
            paramValue:  "abc",
            expectedID:  0,
            expectError: true,
        },
        {
            name:        "Invalid ID - Empty",
            paramValue:  "",
            expectedID:  0,
            expectError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)
            c.Params = []gin.Param{{Key: "id", Value: tt.paramValue}}

            id, err := parseIDParam(c)

            if tt.expectError {
                assert.Error(t, err)
                assert.Equal(t, tt.expectedID, id)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expectedID, id)
            }
        })
    }
}

func TestParsePaginationParams(t *testing.T) {
    gin.SetMode(gin.TestMode)

    tests := []struct {
        name          string
        queryString   string
        expectedPage  int
        expectedLimit int
    }{
        {
            name:          "Default values",
            queryString:   "",
            expectedPage:  1,
            expectedLimit: 100,
        },
        {
            name:          "Custom values",
            queryString:   "page=2&limit=50",
            expectedPage:  2,
            expectedLimit: 50,
        },
        {
            name:          "Invalid values - negative",
            queryString:   "page=-1&limit=-10",
            expectedPage:  1,
            expectedLimit: 100,
        },
        {
            name:          "Invalid values - non-numeric",
            queryString:   "page=abc&limit=def",
            expectedPage:  1,
            expectedLimit: 100,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)
            req := httptest.NewRequest("GET", "/?"+tt.queryString, nil)
            c.Request = req

            page, limit := parsePaginationParams(c)

            assert.Equal(t, tt.expectedPage, page)
            assert.Equal(t, tt.expectedLimit, limit)
        })
    }
}
