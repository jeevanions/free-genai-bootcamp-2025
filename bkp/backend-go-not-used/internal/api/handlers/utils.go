package handlers

import (
    "errors"
    "strconv"

    "github.com/gin-gonic/gin"
)

var (
    ErrInvalidID = errors.New("invalid ID parameter")
)

// parseIDParam extracts and parses an ID parameter from the request URL
func parseIDParam(c *gin.Context) (int64, error) {
    id := c.Param("id")
    if id == "" {
        return 0, ErrInvalidID
    }

    parsedID, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        return 0, ErrInvalidID
    }

    return parsedID, nil
}

// parsePaginationParams extracts page and limit parameters from the request query
func parsePaginationParams(c *gin.Context) (page, limit int) {
    page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
    if page < 1 {
        page = 1
    }

    limit, _ = strconv.Atoi(c.DefaultQuery("limit", "100"))
    if limit < 1 {
        limit = 100
    }

    return page, limit
}
