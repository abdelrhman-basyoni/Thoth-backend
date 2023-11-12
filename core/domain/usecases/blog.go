package domain

import (
	"errors"

	domain "github.com/abdelrhman-basyoni/thoth-backend/core/domain/repositories"
	repos "github.com/abdelrhman-basyoni/thoth-backend/core/implementation/repositories"
	"github.com/abdelrhman-basyoni/thoth-backend/types"
	"gorm.io/gorm"
)

type BlogUseCases struct {
	blogRepo domain.BlogRepository
}

func NewBlogUseCases(db *gorm.DB) *BlogUseCases {
	repo := repos.NewBlogRepoSql(db)
	return &BlogUseCases{blogRepo: repo}
}

func (buc *BlogUseCases) CreateBlog(title, text, authorId string, categories []string) error {

	err := buc.blogRepo.CreateBlog(title, text, authorId, categories)

	if err != nil {
		return errors.New("failed to create Blog")
	}
	return nil

}

func (buc *BlogUseCases) PublishBlog(blogId string, role, authorId string) error {

	if role == types.Roles.Author {
		res := buc.blogRepo.GetBlogForAuthor(blogId, authorId)
		if res == nil {
			return errors.New("invalid Blog")
		}
	}

	err := buc.blogRepo.PublishBlog(blogId)

	if err != nil {
		return errors.New("failed to publish Blog")
	}
	return nil
}

func (buc *BlogUseCases) AddComment(blogId, commenterName, text string) error {
	err := buc.blogRepo.AddComment(blogId, commenterName, text)

	if err != nil {
		return errors.New("failed to Add Comment")
	}
	return nil
}
