package repos

import (
	"strconv"
	"sync"
	"time"

	"github.com/abdelrhman-basyoni/thoth-backend/core/domain/entities"
	"github.com/abdelrhman-basyoni/thoth-backend/core/implementation/models"
	typ "github.com/abdelrhman-basyoni/thoth-backend/types"
	"gorm.io/gorm"
)

type BlogRepoSql struct {
	db *gorm.DB
}

func NewBlogRepoSql(db *gorm.DB) *BlogRepoSql {
	return &BlogRepoSql{db: db}
}

func (br *BlogRepoSql) CreateBlog(title, text, authorId string, categories []string) error {

	authorIdNum, err := strconv.Atoi(authorId)
	if err != nil {
		return err
	}

	res := br.db.Create(&models.Blog{Title: title, Body: text, AuthorId: uint(authorIdNum), Categories: categories})

	return res.Error

}

func (br *BlogRepoSql) PublishBlog(blogId string) error {
	blogNum, err := strconv.Atoi(blogId)
	if err != nil {
		return err
	}
	res := br.db.Model(&models.Blog{}).Where("id = ?", blogNum).Updates(&map[string]interface{}{
		"published":    true,
		"published_at": time.Now(),
	})

	return res.Error
}

func (br *BlogRepoSql) AddComment(blogId, commenterName, text string) error {
	blogNum, err := strconv.ParseUint(blogId, 10, 64)
	if err != nil {
		return err
	}
	res := br.db.Model(&models.Comment{}).Create(&models.Comment{CommenterName: commenterName, Text: text, BlogID: uint(blogNum)})
	return res.Error
}

func (br *BlogRepoSql) GetBlogForAuthor(blogId, authorId string) *entities.Blog {
	authId, err := strconv.ParseUint(blogId, 10, 64)
	if err != nil {
		return nil
	}
	var blog entities.Blog
	res := br.db.First(&blog, "id = ? AND author_id", blogId, uint(authId))

	// Check if a record was found
	if res.RowsAffected == 0 {
		return nil
	}

	return &blog
}

func (br *BlogRepoSql) ApproveComment(commentId string) error {

	commentNum, err := strconv.ParseUint(commentId, 10, 64)
	if err != nil {
		return err
	}
	res := br.db.Model(&models.Comment{}).Where("id = ?", commentNum).UpdateColumn("approved", true)
	return res.Error
}
func (br *BlogRepoSql) DeleteComment(commentId string) error {

	commentNum, err := strconv.ParseUint(commentId, 10, 64)
	if err != nil {
		return err
	}
	res := br.db.Model(&models.Comment{}).Delete("id = ?", commentNum)
	return res.Error
}

func (br *BlogRepoSql) GetBlogComments(blogId string, pageNum int) (*typ.PaginatedEntities[entities.Comment], error) {
	var wg sync.WaitGroup
	wg.Add(2)
	res := typ.PaginatedEntities[entities.Comment]{}
	blogNum, err := strconv.ParseUint(blogId, 10, 64)
	recordsPerPage := 20
	offset := (pageNum - 1) * recordsPerPage

	if err != nil {
		return nil, err
	}
	var totalErr error

	go func() {
		defer wg.Done()
		if err := br.db.Model(&models.Comment{}).Where("blogId = ?", blogNum).Count(&res.Total).Error; err != nil {
			totalErr = err
		}
	}()

	go func() {
		defer wg.Done()
		if err := br.db.Where("blogId = ?", blogNum).Offset(offset).Limit(recordsPerPage).Find(&res.Entities).Error; err != nil {
			totalErr = err
		}
	}()
	// Wait for both goroutines to complete
	wg.Wait()

	if totalErr != nil {
		return nil, totalErr
	}

	return &res, nil
}

func (br *BlogRepoSql) CanUserControlComment(userId, commentId string) (bool, error) {
	commentNum, err := strconv.ParseUint(commentId, 10, 64)
	if err != nil {
		return false, err
	}

	userNum, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		return false, err
	}

	var count int64
	err = br.db.Model(&models.Comment{}).
		Joins("JOIN blogs ON comments.blog_id = blogs.id").
		Where("comments.id = ? AND blogs.author_id = ?", commentNum, userNum).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil

}
