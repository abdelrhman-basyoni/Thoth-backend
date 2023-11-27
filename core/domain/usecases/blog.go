package domain

import (
	"errors"

	"github.com/abdelrhman-basyoni/thoth-backend/core/domain/entities"
	domain "github.com/abdelrhman-basyoni/thoth-backend/core/domain/repositories"
	repos "github.com/abdelrhman-basyoni/thoth-backend/core/implementation/repositories"
	typ "github.com/abdelrhman-basyoni/thoth-backend/types"

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

	if role == typ.Roles.Author {
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

func (buc *BlogUseCases) GetPublishedBlogById(blogId string) (*entities.Blog, error) {
	res := buc.blogRepo.GetBlogById(blogId, true)

	if res == nil {
		return nil, errors.New("invalid blog id")
	}

	return res, nil
}

func (buc *BlogUseCases) ApproveComment(commentId, userId, role string) error {

	if role != typ.Roles.Admin {
		check, _ := buc.blogRepo.CanUserControlBlog(userId, commentId)
		if !check {
			return errors.New("unauthorized to Approve Comment")
		}
	}

	err := buc.blogRepo.ApproveComment(commentId)

	if err != nil {
		return errors.New("failed to Approve Comment")
	}
	return nil

}

func (buc *BlogUseCases) DeleteComment(commentId, userId, role string) error {
	if role != typ.Roles.Admin {
		check, _ := buc.blogRepo.CanUserControlBlog(userId, commentId)
		if !check {
			return errors.New("unauthorized to Approve Comment")
		}
	}

	err := buc.blogRepo.DeleteComment(commentId)

	if err != nil {
		return errors.New("failed to Approve Comment")
	}
	return nil

}

func (buc *BlogUseCases) GetBlogComments(blogId string, pageNum int) (*typ.PaginatedEntities[entities.Comment], error) {
	if pageNum <= 0 {
		return nil, errors.New("invalid page number")
	}
	return buc.blogRepo.GetBlogComments(blogId, pageNum)
}
