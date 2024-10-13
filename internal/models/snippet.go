package models

import (
	"time"
)

// Snippet represents a snippet persisted to storage.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
