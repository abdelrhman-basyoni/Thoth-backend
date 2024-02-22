package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Blog struct {
	gorm.Model
	Title       string         `json:"title"`
	Body        string         `json:"body" gorm:"type:text"`
	Categories  pq.StringArray `json:"categories" gorm:"type:text[]"`
	Image       string         `json:"image"`
	AuthorId    uint           `json:"authorId" gorm:"foreignKey:UserID"`
	Published   bool           `json:"published"`
	PublishedAt time.Time      `json:"publishedAt"`
	Comments    []Comment
}

type Comment struct {
	gorm.Model
	BlogID        uint   `json:"blogId"`
	CommenterName string `json:"commenter"`
	Text          string `json:"text"`
	Approved      bool   `json:"approved"`
}
