package domain

import (
	"time"

	"github.com/abdelrhman-basyoni/thoth-backend/core/domain/entities"
	typ "github.com/abdelrhman-basyoni/thoth-backend/types"
)

type BlogAuthorData struct {
	ID   uint   `json:"id" gorm:"embedded;embeddedPrefix:author_"`
	Name string `json:"name" gorm:"embedded;embeddedPrefix:author_"`
}
type BlogData struct {
	ID          uint           `json:"id"`
	Title       string         `json:"title"`
	Body        string         `json:"body"`
	Published   bool           `json:"published"`
	PublishedAt time.Time      `json:"publishedAt"`
	Categories  []string       `json:"categories"`
	Author      BlogAuthorData `json:"author"`
}

type BlogRepository interface {
	CreateBlog(title, text string, authorId uint, categories []string) error
	PublishBlog(blogId uint) error
	GetBlogsFiltered(authorId *uint, category *string, pageNum int) (*typ.PaginatedEntities[BlogData], error)
	GetBlogById(blogId uint, mustBePublished bool) *BlogData
	AddComment(blogId uint, commenterName, text string) error
	GetBlogForAuthor(blogId, authorId uint) *entities.Blog
	GetBlogComments(blogId uint, pageNum int) (*typ.PaginatedEntities[entities.Comment], error)
	ApproveComment(commentId uint) error
	DeleteComment(commentId uint) error
	CanUserControlBlog(userId, commentId uint) (bool, error)
}
