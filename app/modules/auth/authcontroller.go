package authModule

import (
	"net/http"

	domain "github.com/abdelrhman-basyoni/thoth-backend/core/domain/usecases"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var validate = validator.New()

type createUser struct {
	Name     string `json:"name" validate:"required" `
	Email    string `json:"email" validate:"required,email" `
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required"`
}

type loginDto struct {
	Email    string `json:"email" validate:"required,email" `
	Password string `json:"password" validate:"required"`
}

type AuthController struct {
	uc *domain.AuthUseCases
}

func NewAuthController(db *gorm.DB) *AuthController {
	useCases := domain.NewAuthUseCases(db)

	return &AuthController{uc: useCases}
}

func (ac *AuthController) HandleCreate(c echo.Context) error {

	var user createUser

	// Bind and validate the request body
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
		})
	}

	// use the validator library to validate required fields
	if err := validate.Struct(&user); err != nil {

		return err
	}

	if err := ac.uc.Create(user.Name, user.Email, user.Password, user.Role); err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)

}

func (ac *AuthController) HandleLogin(c echo.Context) error {
	var user loginDto

	// Bind and validate the request body
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
		})
	}

	// use the validator library to validate required fields
	if err := validate.Struct(&user); err != nil {

		return err
	}
	token, err := ac.uc.Login(user.Email, user.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (ac *AuthController) Test(c echo.Context) error {
	return c.NoContent(http.StatusCreated)
}
