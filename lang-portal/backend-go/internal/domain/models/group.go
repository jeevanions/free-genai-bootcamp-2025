package models

type GroupResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	WordCount int    `json:"word_count"`
}

type GroupListResponse struct {
	Items      []GroupResponse    `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}

type GroupDetailResponse struct {
	ID        int64 `json:"id"`
	Name      string `json:"name"`
	Stats     GroupStats `json:"stats"`
}

type GroupStats struct {
	TotalWordCount int `json:"total_word_count"`
}

type GroupWordsResponse struct {
	Items      []WordResponse     `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}

type GroupStudySessionsResponse struct {
	Items      []StudySessionResponse `json:"items"`
	Pagination PaginationResponse     `json:"pagination"`
}
