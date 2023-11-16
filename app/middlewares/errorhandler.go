package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func GlobalErrorHandler(err error, c echo.Context) {
	// Log the error or perform any other custom error handling logic
	log.Error(err)

	// Respond with an appropriate error message and status code
	c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"error": err.Error(),
	})
}
