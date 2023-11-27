package domain

import (
	"github.com/abdelrhman-basyoni/thoth-backend/core/domain/entities"
	typ "github.com/abdelrhman-basyoni/thoth-backend/types"
)

type BlogRepository interface {
	CreateBlog(title, text, authorId string, categories []string) error
	PublishBlog(blogId string) error
	GetBlogById(blogId string, mustBePublished bool) *entities.Blog
	AddComment(blogId, commenterName, text string) error
	GetBlogForAuthor(blogId, authorId string) *entities.Blog
	GetBlogComments(blogId string, pageNum int) (*typ.PaginatedEntities[entities.Comment], error)
	ApproveComment(commentId string) error
	DeleteComment(commentId string) error
	CanUserControlBlog(userId, commentId string) (bool, error)
}
