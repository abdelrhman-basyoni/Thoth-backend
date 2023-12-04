package userModule

import (
	"github.com/abdelrhman-basyoni/thoth-backend/app/middlewares"
	typ "github.com/abdelrhman-basyoni/thoth-backend/types"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterUserRoutes(e *echo.Echo, db *gorm.DB) {
	controller := NewUserController(db)
	authGroup := e.Group("/user")
	authGroup.GET("/all", controller.HandleGetAllUsers, middlewares.RoleAuth([]string{typ.Roles.Admin}))

}
