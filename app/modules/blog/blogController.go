package blogModule

import (
	"net/http"
	"reflect"
	"strconv"

	domain "github.com/abdelrhman-basyoni/thoth-backend/core/domain/usecases"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var validate = validator.New()

type createBlog struct {
	Title      string   `json:"title" validate:"required" `
	Text       string   `json:"text" validate:"required" `
	Categories []string `json:"categories" validate:"required,isStringArray"`
}

func isStringArray(fl validator.FieldLevel) bool {
	// Check if the field is a slice or array
	if fl.Field().Kind() != reflect.Slice && fl.Field().Kind() != reflect.Array {
		return false
	}

	// Check if each element is a string
	for i := 0; i < fl.Field().Len(); i++ {
		if fl.Field().Index(i).Kind() != reflect.String {
			return false
		}
	}

	return true
}

type BlogController struct {
	uc *domain.BlogUseCases
}

func NewBlogController(db *gorm.DB) *BlogController {
	useCases := domain.NewBlogUseCases(db)

	return &BlogController{uc: useCases}
}

func (bc *BlogController) HandleCreate(c echo.Context) error {

	var blog createBlog
	userId := c.Get("user").(string)
	// Bind and validate the request body
	if err := c.Bind(&blog); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
		})
	}
	if err := validate.RegisterValidation("isStringArray", isStringArray); err != nil {
		return err
	}

	// use the validator library to validate required fields
	if err := validate.Struct(&blog); err != nil {

		return err
	}

	if err := bc.uc.CreateBlog(blog.Title, blog.Text, userId, blog.Categories); err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)

}

func (bc *BlogController) HandlePublish(c echo.Context) error {
	role := c.Get("userRole").(string)
	userId := c.Get("user").(string)
	blogID := c.Param("id")

	if err := bc.uc.PublishBlog(blogID, role, userId); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

type addComment struct {
	CommenterName string `json:"commenterName" validate:"required" `
	Text          string `json:"text" validate:"required" `
}

func (bc *BlogController) HandleAddComment(c echo.Context) error {
	blogID := c.Param("id")
	var comment addComment
	// Bind and validate the request body
	if err := c.Bind(&comment); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
		})
	}

	// use the validator library to validate required fields
	if err := validate.Struct(&comment); err != nil {

		return err
	}

	if err := bc.uc.AddComment(blogID, comment.CommenterName, comment.Text); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (bc *BlogController) HandleGetPublishedBlog(c echo.Context) error {
	blogId := c.Param("id")

	res, err := bc.uc.GetPublishedBlogById(blogId)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"blog": res,
	})
}

func (bc *BlogController) HandleApproveComment(c echo.Context) error {
	userId := c.Get("user").(string)
	userRole := c.Get("role").(string)
	commentId := c.Param("id")
	err := bc.uc.ApproveComment(commentId, userId, userRole)

	if err.Error() == "unauthorized to Approve Comment" {
		return c.NoContent(http.StatusForbidden)
	}

	return c.NoContent(http.StatusOK)
}

func (bc *BlogController) HandleDeleteComment(c echo.Context) error {
	commentId := c.Param("id")
	userId := c.Get("user").(string)
	userRole := c.Get("role").(string)
	if err := bc.uc.DeleteComment(commentId, userId, userRole); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid comment id",
		})
	}
	return c.NoContent(http.StatusOK)
}

func (bc *BlogController) HandleGetBlogComments(c echo.Context) error {
	blogId := c.Param("id")
	pageParam := c.QueryParam("page")
	pageNum := 1
	if pageParam != "" {
		pageVal, err := strconv.ParseInt(pageParam, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "Invalid page number",
			})
		}
		pageNum = int(pageVal)
	}

	res, err := bc.uc.GetBlogComments(blogId, pageNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)

}

func (bc *BlogController) HandleGetBlogs(c echo.Context) error {
	pageParam := c.QueryParam("page")
	authorId := c.QueryParam("author")
	category := c.QueryParam("category")
	pageNum := 1

	if pageParam != "" {
		pageVal, err := strconv.ParseInt(pageParam, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "Invalid page number",
			})
		}
		pageNum = int(pageVal)
	}
	res, err := bc.uc.GetAllBlogsPaginated(&authorId, &category, pageNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"blogs": res,
	})
}
