package userModule

import (
	"net/http"

	domain "github.com/abdelrhman-basyoni/thoth-backend/core/domain/usecases"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserController struct {
	uc *domain.UserUseCases
}

func NewUserController(db *gorm.DB) *UserController {
	useCases := domain.NewUserUseCases(db)

	return &UserController{uc: useCases}
}

func (usc *UserController) HandleGetAllUsers(c echo.Context) error {

	res, err := usc.uc.GetAllUsers()

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
