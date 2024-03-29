package blogModule

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	domain "github.com/abdelrhman-basyoni/thoth-backend/core/domain/usecases"
	typ "github.com/abdelrhman-basyoni/thoth-backend/types"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

var validate = validator.New()

type createBlog struct {
	Title      string   `json:"title" validate:"required" `
	Text       string   `json:"body" validate:"required" `
	Image      string   `json:"image" validate:"required"`
	Publish    bool     `json:"publish"`
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
	userId := c.Get("user").(uint)
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

	if err := bc.uc.CreateBlog(blog.Title, blog.Text, userId, blog.Categories, blog.Publish || false, blog.Image); err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)

}

func (bc *BlogController) HandlePublishToggle(c echo.Context) error {
	role := c.Get("userRole").(string)
	userId := c.Get("user").(uint)
	blogID := c.Param("id")
	publish := c.QueryParam("publish")

	blogIdUint, err := strconv.ParseUint(blogID, 10, 64)
	if err != nil {
		// Handle the error if the conversion fails
		return c.JSON(http.StatusBadRequest, "Invalid blog ID")

	}
	var pub = true
	if publish == "false" {
		pub = false
	}
	if err := bc.uc.TogglePublishBlog(uint(blogIdUint), userId, role, pub); err != nil {
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
	blogIdUint, err := strconv.ParseUint(blogID, 10, 64)
	if err != nil {
		// Handle the error if the conversion fails
		return c.JSON(http.StatusBadRequest, "Invalid blog ID")

	}
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

	if err := bc.uc.AddComment(uint(blogIdUint), comment.CommenterName, comment.Text); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (bc *BlogController) HandleGetPublishedBlog(c echo.Context) error {
	blogId := c.Param("id")
	blogIdUint, err := strconv.ParseUint(blogId, 10, 64)
	if err != nil {
		// Handle the error if the conversion fails
		return c.JSON(http.StatusBadRequest, "Invalid blog ID")

	}
	res, err := bc.uc.GetPublishedBlogById(uint(blogIdUint))
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (bc *BlogController) HandleGetMyBlog(c echo.Context) error {
	blogId := c.Param("id")
	fmt.Println(blogId)
	blogIdUint, err := strconv.ParseUint(blogId, 10, 64)
	if err != nil {
		// Handle the error if the conversion fails
		return c.JSON(http.StatusBadRequest, "Invalid blog ID")

	}
	res, err := bc.uc.GetMyBlogById(uint(blogIdUint))
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (bc *BlogController) HandleApproveComment(c echo.Context) error {
	userId := c.Get("user").(uint)
	userRole := c.Get("userRole").(string)
	commentId := c.Param("id")

	// Convert commentId to a uint
	commentIdUint, err := strconv.ParseUint(commentId, 10, 64)
	if err != nil {
		// Handle the error if the conversion fails
		return c.JSON(http.StatusBadRequest, "Invalid comment ID")

	}

	err = bc.uc.ApproveComment(uint(commentIdUint), userId, userRole)
	if err != nil {
		if err.Error() == "unauthorized to Approve Comment" {
			return c.NoContent(http.StatusForbidden)
		}
	}

	return c.NoContent(http.StatusOK)
}

func (bc *BlogController) HandleDeleteComment(c echo.Context) error {
	commentId := c.Param("id")
	commentIdUint, err := strconv.ParseUint(commentId, 10, 64)
	if err != nil {
		// Handle the error if the conversion fails
		return c.JSON(http.StatusBadRequest, "Invalid comment ID")

	}
	userId := c.Get("user").(uint)
	userRole := c.Get("userRole").(string)
	if err := bc.uc.DeleteComment(uint(commentIdUint), userId, userRole); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid comment id",
		})
	}
	return c.NoContent(http.StatusOK)
}

func (bc *BlogController) HandleGetBlogComments(c echo.Context) error {
	blogId := c.Param("id")
	blogIdUint, err := strconv.ParseUint(blogId, 10, 64)
	if err != nil {
		// Handle the error if the conversion fails
		return c.JSON(http.StatusBadRequest, "Invalid blog ID")

	}
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

	res, err := bc.uc.GetBlogComments(uint(blogIdUint), pageNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)

}

func (bc *BlogController) HandleGetMyBlogComments(c echo.Context) error {
	blogId := c.Param("id")
	blogIdUint, err := strconv.ParseUint(blogId, 10, 64)
	if err != nil {
		// Handle the error if the conversion fails
		return c.JSON(http.StatusBadRequest, "Invalid blog ID")

	}
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

	res, err := bc.uc.GetMyBlogComments(uint(blogIdUint), pageNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)

}

func (bc *BlogController) HandleGetBlogNotApprovedComments(c echo.Context) error {
	blogId := c.Param("id")
	blogIdUint, err := strconv.ParseUint(blogId, 10, 64)
	if err != nil {
		// Handle the error if the conversion fails
		return c.JSON(http.StatusBadRequest, "Invalid blog ID")

	}
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

	res, err := bc.uc.GetBlogNotApprovedComments(uint(blogIdUint), pageNum)
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
	authorIdUint, err := strconv.ParseUint(authorId, 10, 64)
	authorU := uint(authorIdUint)
	authID := &authorU
	if err != nil {
		authID = nil

	}
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
	res, err := bc.uc.GetAllBlogsPaginated(authID, &category, pageNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"blogs": res,
	})
}

func (bc *BlogController) HandleGetMyBlogs(c echo.Context) error {
	role := c.Get("userRole").(string)
	userId := c.Get("user").(uint)
	authorId := &userId

	if role == typ.Roles.Admin {
		authorId = nil
	}

	res, err := bc.uc.GetAllMyBlogsPaginated(authorId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

type editBlog struct {
	Title string `json:"title" validate:"required" `
	Body  string `json:"body" validate:"required" `
}

func (bc *BlogController) HandleEditBlog(c echo.Context) error {
	role := c.Get("userRole").(string)
	userId := c.Get("user").(uint)
	blogId := c.Param("id")

	blogIdUint, err := strconv.ParseUint(blogId, 10, 64)
	if err != nil {
		// Handle the error if the conversion fails
		return c.JSON(http.StatusBadRequest, "Invalid blog ID")

	}

	var blogContent editBlog
	// Bind and validate the request body
	if err := c.Bind(&blogContent); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
		})
	}

	err = bc.uc.EditBlog(role, userId, uint(blogIdUint), blogContent.Title, blogContent.Body)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)

}
