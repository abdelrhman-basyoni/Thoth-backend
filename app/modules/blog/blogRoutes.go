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
	BlogGroup.POST("/create", controller.HandleCreate, middlewares.RoleAuth([]string{typ.Roles.Admin, typ.Roles.Author}))
	BlogGroup.GET("/getAll", controller.HandleGetBlogs)
	BlogGroup.GET("/getMyBlogs", controller.HandleGetMyBlogs, middlewares.RoleAuth([]string{typ.Roles.Admin, typ.Roles.Author}))
	BlogGroup.PUT("/edit/:id", controller.HandleEditBlog, middlewares.RoleAuth([]string{typ.Roles.Admin, typ.Roles.Author}))
	BlogGroup.GET("/get/:id", controller.HandleGetPublishedBlog)
	BlogGroup.GET("/getMyBlog/:id", controller.HandleGetMyBlog, middlewares.RoleAuth([]string{typ.Roles.Admin, typ.Roles.Author}))
	BlogGroup.PUT("/togglePublish/:id", controller.HandlePublishToggle, middlewares.RoleAuth([]string{typ.Roles.Admin, typ.Roles.Author}))
	BlogGroup.GET("/comments/approved/:id", controller.HandleGetBlogComments)
	BlogGroup.GET("/comments/my/:id", controller.HandleGetMyBlogComments, middlewares.RoleAuth([]string{typ.Roles.Admin, typ.Roles.Author}))
	BlogGroup.GET("/comments/notApproved/:id", controller.HandleGetBlogNotApprovedComments, middlewares.RoleAuth([]string{typ.Roles.Admin, typ.Roles.Author}))
	BlogGroup.PUT("/comments/approve/:id", controller.HandleApproveComment, middlewares.RoleAuth([]string{typ.Roles.Admin, typ.Roles.Author}))
	BlogGroup.POST("/comments/:id", controller.HandleAddComment)
	BlogGroup.DELETE("/comments/delete/:id", controller.HandleDeleteComment, middlewares.RoleAuth([]string{typ.Roles.Admin, typ.Roles.Author}))

}
