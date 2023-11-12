package domain

import "github.com/abdelrhman-basyoni/thoth-backend/core/domain/entities"

type BlogRepository interface {
	CreateBlog(title, text, authorId string, categories []string) error
	PublishBlog(blogId string) error
	AddComment(blogId, commenterName, text string) error
	GetBlogForAuthor(blogId, authorId string) *entities.Blog
}
