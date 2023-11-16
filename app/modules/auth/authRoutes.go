package authModule

import (
	"github.com/abdelrhman-basyoni/thoth-backend/app/middlewares"
	"github.com/abdelrhman-basyoni/thoth-backend/types"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(e *echo.Echo, db *gorm.DB) {
	controller := NewAuthController(db)
	authGroup := e.Group("/auth")
	authGroup.POST("/create", controller.HandleCreate, middlewares.RoleAuth([]string{types.Roles.Admin}))
	authGroup.POST("/login", controller.HandleLogin)
	authGroup.GET("/test", controller.Test, middlewares.RoleAuth([]string{types.Roles.Admin}))

}
