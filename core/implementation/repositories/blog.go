package repos

import (
	"strconv"
	"time"

	"github.com/abdelrhman-basyoni/thoth-backend/core/domain/entities"
	"github.com/abdelrhman-basyoni/thoth-backend/core/implementation/models"
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
