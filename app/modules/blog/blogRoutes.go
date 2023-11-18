package blogModule

import (
	"github.com/abdelrhman-basyoni/thoth-backend/app/middlewares"
	typ "github.com/abdelrhman-basyoni/thoth-backend/types"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterBlogRoutes(e *echo.Echo, db *gorm.DB) {
	controller := NewBlogController(db)
	BlogGroup := e.Group("/blog")
	BlogGroup.POST("/create", controller.HandleCreate, middlewares.RoleAuth([]string{typ.Roles.Admin}))
	BlogGroup.POST("/publish/:id", controller.HandlePublish, middlewares.RoleAuth([]string{typ.Roles.Admin}))
	BlogGroup.POST("/comments/:id", controller.HandleAddComment)

}
