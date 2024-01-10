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

type CommentData struct {
	ID            uint   `json:"id"`
	BlogID        uint   `json:"blogId"`
	CommenterName string `json:"commenterName"`
	Text          string `json:"text"`
	Approved      bool   `json:"approved"`
}

type BlogRepository interface {
	CreateBlog(title, text string, authorId uint, categories []string, publish bool) error
	EditBlog(blogId uint, title, body string) error
	TogglePublishBlog(blogId uint, publish bool) error
	GetBlogsFiltered(authorId *uint, category *string, pageNum int) (*typ.PaginatedEntities[BlogData], error)
	GetBlogById(blogId uint, mustBePublished bool) *BlogData
	AddComment(blogId uint, commenterName, text string) error
	GetBlogForAuthor(blogId, authorId uint) *entities.Blog
	GetBlogComments(blogId uint, pageNum int) (*typ.PaginatedEntities[entities.Comment], error)
	GetMyBlogComments(blogId uint, pageNum int) (*typ.PaginatedEntities[CommentData], error)
	GetBlogNotApprovedComments(blogId uint, pageNum int) (*typ.PaginatedEntities[entities.Comment], error)
	ApproveComment(commentId uint) error
	DeleteComment(commentId uint) error
	CanUserControlBlog(userId, commentId uint) (bool, error)
	CanUserControlComment(userId, commentId uint) (bool, error)
}
