package models

import "time"

// Valid values for group categories
var ValidGroupCategories = []string{"grammar", "thematic", "situational"}

// Group represents a thematic group in the database
type Group struct {
	ID              int64     `db:"id" json:"id"`
	Name            string    `db:"name" json:"name"`
	Description     string    `db:"description" json:"description"`
	DifficultyLevel int       `db:"difficulty_level" json:"difficulty_level"`
	Category        string    `db:"category" json:"category"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

// TableName returns the name of the table for the Group model
func (Group) TableName() string {
	return "groups"
}

// WordGroup represents the many-to-many relationship between words and groups
type WordGroup struct {
	ID        int64     `db:"id" json:"id"`
	WordID    int64     `db:"word_id" json:"word_id"`
	GroupID   int64     `db:"group_id" json:"group_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// TableName returns the name of the table for the WordGroup model
func (WordGroup) TableName() string {
	return "words_groups"
}
