package repos

import (
	"fmt"
	"sync"
	"time"

	"github.com/abdelrhman-basyoni/thoth-backend/app/config"
	"github.com/abdelrhman-basyoni/thoth-backend/core/domain/entities"
	domain "github.com/abdelrhman-basyoni/thoth-backend/core/domain/repositories"
	"github.com/abdelrhman-basyoni/thoth-backend/core/implementation/models"
	typ "github.com/abdelrhman-basyoni/thoth-backend/types"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type BlogData struct {
	ID          uint           `json:"id"`
	Title       string         `json:"title"`
	Body        string         `json:"body"`
	Published   bool           `json:"published"`
	PublishedAt time.Time      `json:"publishedAt"`
	Categories  pq.StringArray `json:"categories" gorm:"type:text[]"`
	AuthorName  string         `json:"author_name"`
	AuthorID    uint           `json:"author_id"`
}

func (res *BlogData) fmtToDomainBlogData() domain.BlogData {
	blog := domain.BlogData{ID: res.ID, Title: res.Title, Body: res.Body, Published: res.Published, Author: domain.BlogAuthorData{Name: res.AuthorName, ID: res.AuthorID}, PublishedAt: res.PublishedAt, Categories: res.Categories}
	return blog
}

type BlogRepoSql struct {
	db *gorm.DB
}

func NewBlogRepoSql(db *gorm.DB) *BlogRepoSql {
	return &BlogRepoSql{db: db}
}

func (br *BlogRepoSql) CreateBlog(title, text string, authorId uint, categories []string, publish bool) error {
	var res *gorm.DB
	if publish {
		res = br.db.Create(&models.Blog{Title: title, Body: text, AuthorId: authorId, Categories: categories, Published: publish, PublishedAt: time.Now()})

	} else {
		res = br.db.Create(&models.Blog{Title: title, Body: text, AuthorId: authorId, Categories: categories})

	}

	return res.Error

}

func (br *BlogRepoSql) GetBlogById(blogId uint, mustBePublished bool) *domain.BlogData {
	selectQuery := "blogs.id, blogs.title, blogs.body, blogs.published, blogs.published_at, blogs.categories, users.id as author_id,  users.name as author_name"
	var blog BlogData

	if mustBePublished {
		res := br.db.Model(&models.Blog{}).Select(selectQuery).
			Joins("JOIN users ON blogs.author_id = users.id").First(&blog, "blogs.id = ? AND published = true", blogId).Limit(1)

		// Check if a record was found
		if res.RowsAffected == 0 {
			return nil
		}
		fmt.Println(blog)
		fmtRes := blog.fmtToDomainBlogData()
		fmt.Println(fmtRes)
		return &fmtRes
	}

	res := br.db.Model(&models.Blog{}).Select(selectQuery).
		Joins("JOIN users ON blogs.author_id = users.id").First(&blog, "blogs.id = ?", blogId)

	// Check if a record was found
	if res.RowsAffected == 0 {
		return nil
	}
	fmtRes := blog.fmtToDomainBlogData()

	return &fmtRes

}

func (br *BlogRepoSql) GetBlogsFiltered(authorId *uint, category *string, pageNum int) (*typ.PaginatedEntities[domain.BlogData], error) {

	validAuthor := authorId != nil
	validCategory := category != nil && *category != ""
	res := typ.PaginatedEntities[domain.BlogData]{}
	res1 := []BlogData{}
	test := &category
	var whereQuery = "published = true"
	var variables []any

	if validAuthor && validCategory {

		whereQuery += " AND author_id = ? AND ? = ANY(categories) "
		variables = append(variables, authorId, test)
	} else if validAuthor && !validCategory {

		whereQuery += " AND author_id = ?"
		variables = append(variables, authorId)
	} else if !validAuthor && validCategory {
		whereQuery += " AND ? = ANY(categories)"
		variables = append(variables, test)
	}
	var wg sync.WaitGroup
	wg.Add(2)

	offset := (pageNum - 1) * config.RecordsPerPage
	var totalErr error
	selectQuery := "blogs.id, blogs.title, blogs.body, blogs.published, blogs.published_at, blogs.categories, users.id as author_id,  users.name as author_name"
	go func() {
		defer wg.Done()
		if err := br.db.Model(&models.Blog{}).Select(selectQuery).
			Joins("JOIN users ON blogs.author_id = users.id").
			Where(whereQuery, variables...).Count(&res.Total).Error; err != nil {
			totalErr = err
		}
	}()

	go func() {
		defer wg.Done()
		if err := br.db.Model(&models.Blog{}).Select(selectQuery).
			Joins("JOIN users ON blogs.author_id = users.id").
			Where(whereQuery, variables...).Offset(offset).Limit(config.RecordsPerPage).Find(&res1).Error; err != nil {
			totalErr = err
		}
	}()
	// Wait for both goroutines to complete
	wg.Wait()

	if totalErr != nil {
		return nil, totalErr
	}
	entities := res.Entities
	for _, res := range res1 {
		blog := domain.BlogData{ID: res.ID, Title: res.Title, Body: res.Body, Published: res.Published, Author: domain.BlogAuthorData{Name: res.AuthorName, ID: res.AuthorID}, PublishedAt: res.PublishedAt, Categories: res.Categories}
		entities = append(entities, blog)
	}
	res.Entities = entities

	return &res, nil

}

func (br *BlogRepoSql) TogglePublishBlog(blogId uint, publish bool) error {

	res := br.db.Model(&models.Blog{}).Where("id = ?", blogId).Updates(&map[string]interface{}{
		"published":    publish,
		"published_at": time.Now(),
	})

	return res.Error
}

func (br *BlogRepoSql) AddComment(blogId uint, commenterName, text string) error {

	res := br.db.Model(&models.Comment{}).Create(&models.Comment{CommenterName: commenterName, Text: text, BlogID: blogId})
	return res.Error
}

func (br *BlogRepoSql) GetBlogForAuthor(blogId, authorId uint) *entities.Blog {

	var blog entities.Blog
	res := br.db.First(&blog, "id = ? AND author_id", blogId, authorId)

	// Check if a record was found
	if res.RowsAffected == 0 {
		return nil
	}

	return &blog
}

func (br *BlogRepoSql) ApproveComment(commentId uint) error {

	res := br.db.Model(&models.Comment{}).Where("id = ?", commentId).UpdateColumn("approved", true)

	return res.Error
}

func (br *BlogRepoSql) DeleteComment(commentId uint) error {

	res := br.db.Model(&models.Comment{}).Delete("id = ?", commentId)
	return res.Error
}

func (br *BlogRepoSql) GetBlogComments(blogId uint, pageNum int) (*typ.PaginatedEntities[entities.Comment], error) {
	var wg sync.WaitGroup
	wg.Add(2)
	res := typ.PaginatedEntities[entities.Comment]{}

	offset := (pageNum - 1) * config.RecordsPerPage

	var totalErr error

	go func() {
		defer wg.Done()
		if err := br.db.Model(&models.Comment{}).Where("blog_id = ? AND approved = true", blogId).Count(&res.Total).Error; err != nil {
			totalErr = err
		}
	}()

	go func() {
		defer wg.Done()
		if err := br.db.Model(&models.Comment{}).Where("blog_id = ?", blogId).Offset(offset).Limit(config.RecordsPerPage).Find(&res.Entities).Error; err != nil {
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

func (br *BlogRepoSql) GetMyBlogComments(blogId uint, pageNum int) (*typ.PaginatedEntities[domain.CommentData], error) {
	var wg sync.WaitGroup
	wg.Add(2)
	res := typ.PaginatedEntities[domain.CommentData]{}

	offset := (pageNum - 1) * config.RecordsPerPage

	var totalErr error

	go func() {
		defer wg.Done()
		if err := br.db.Model(&models.Comment{}).Where("blog_id = ?", blogId).Count(&res.Total).Error; err != nil {
			totalErr = err
		}
	}()

	go func() {
		defer wg.Done()
		if err := br.db.Model(&models.Comment{}).Where("blog_id = ?", blogId).Offset(offset).Limit(config.RecordsPerPage).Find(&res.Entities).Error; err != nil {
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

func (br *BlogRepoSql) GetBlogNotApprovedComments(blogId uint, pageNum int) (*typ.PaginatedEntities[entities.Comment], error) {
	var wg sync.WaitGroup
	wg.Add(2)
	res := typ.PaginatedEntities[entities.Comment]{}

	offset := (pageNum - 1) * config.RecordsPerPage

	var totalErr error

	go func() {
		defer wg.Done()
		if err := br.db.Model(&models.Comment{}).Where("blogId = ? AND approved = false", blogId).Count(&res.Total).Error; err != nil {
			totalErr = err
		}
	}()

	go func() {
		defer wg.Done()
		if err := br.db.Model(&models.Comment{}).Where("blogId = ?", blogId).Offset(offset).Limit(config.RecordsPerPage).Find(&res.Entities).Error; err != nil {
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

func (br *BlogRepoSql) CanUserControlComment(userId, commentId uint) (bool, error) {

	var count int64
	err := br.db.Model(&models.Comment{}).
		Joins("JOIN blogs ON comments.blog_id = blogs.id").
		Where("comments.id = ? AND blogs.author_id = ?", commentId, userId).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil

}

func (br *BlogRepoSql) CanUserControlBlog(userId, blogId uint) (bool, error) {

	var count int64
	err := br.db.Model(&models.Blog{}).
		Where("id = ? AND author_id = ?", blogId, userId).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil

}

func (br *BlogRepoSql) EditBlog(blogId uint, title, body string) error {
	res := br.db.Model(&models.Blog{}).Where("id = ?", blogId).Updates(&map[string]interface{}{
		"title": title,
		"body":  body,
	})

	return res.Error

}
