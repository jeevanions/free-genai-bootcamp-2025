package models

type WordResponse struct {
	ID           int64                  `json:"id"`
	Italian      string                 `json:"italian"`
	English      string                 `json:"english"`
	Parts        map[string]interface{} `json:"parts"`
	CorrectCount int                    `json:"correct_count"`
	WrongCount   int                    `json:"wrong_count"`
}

type WordListResponse struct {
	Items      []WordResponse `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}


