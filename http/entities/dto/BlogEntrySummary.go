package dto

import "time"

// BlogEntry represents a row in the blog_entries table.
type BlogEntrySummary struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Author    string    `json:"author" db:"author"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
