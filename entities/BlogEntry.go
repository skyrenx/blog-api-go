package entities

import "time"

// BlogEntry represents a row in the blog_entries table.
type BlogEntry struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	Author    string    `json:"author" db:"author"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Published bool      `json:"published" db:"published"`
}
