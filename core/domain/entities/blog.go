package entities

import "time"

type Blog struct {
	ID          uint      `json:"id"` // defining id as string so it can work with any database not just sql types
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	Published   bool      `json:"published"`
	PublishedAt time.Time `json:"publishedAt"`
	Categories  []string  `json:"categories"`
	Comments    []Comment `json:"comments"`
}

type Comment struct {
	CommenterName string `json:"commenter"`
	Text          string `json:"text"`
	Approved      bool   `json:"approved"`
}
