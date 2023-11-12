package middlewares

import (
	"fmt"
	"net/http"

	"github.com/abdelrhman-basyoni/thoth-backend/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func RoleAuth(roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the token from the request
			tokenString := c.Request().Header.Get("Authorization")

			// Validate and parse the token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Verify the signing method and return the secret key
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("invalid signing method")
				}

				return []byte(utils.ReadEnv("SECRET_KEY")), nil
			})

			if err != nil || !token.Valid {
				fmt.Println(err.Error())
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}

			// // Access the claims from the token
			claims, ok1 := token.Claims.(jwt.MapClaims)
			userId, ok2 := claims["id"].(string)

			if !ok1 || !ok2 {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}

			tokenRole, ok := claims["role"].(string)

			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid role")
			}

			if !utils.ContainsValue(roles, tokenRole) {
				// Authorized access for admin
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid role")

			}

			// Set user data in the context for use in handlers

			c.Set("user", userId)
			c.Set("userRole", tokenRole)

			return next(c)
		}
	}
}
