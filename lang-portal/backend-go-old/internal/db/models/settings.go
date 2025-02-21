package models

type Settings struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Theme  string `json:"theme"`
	Language string `json:"language"`
}

type Preferences struct {
	ID     int64                  `json:"id"`
	UserID int64                  `json:"user_id"`
	Data   map[string]interface{} `json:"data"`
}
