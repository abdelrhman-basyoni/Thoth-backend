package main

import (
	"fmt"

	"github.com/abdelrhman-basyoni/thoth-backend/app/middlewares"
	authModule "github.com/abdelrhman-basyoni/thoth-backend/app/modules/auth"
	blogModule "github.com/abdelrhman-basyoni/thoth-backend/app/modules/blog"
	userModule "github.com/abdelrhman-basyoni/thoth-backend/app/modules/user"
	"github.com/abdelrhman-basyoni/thoth-backend/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	const port = 3000
	db := config.InitDB()
	e := echo.New()
	e.Use(middleware.CORS())
	db.Debug()
	// Set the router to be case-insensitive

	blogModule.RegisterBlogRoutes(e, db)
	authModule.RegisterAuthRoutes(e, db)
	userModule.RegisterUserRoutes(e, db)
	// Middleware
	e.HTTPErrorHandler = middlewares.GlobalErrorHandler
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
