package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// simple jwt authentication middleware
func AuthenticationMiddleware(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, "error on fetching authorization header")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader || tokenString == "" {
			return c.JSON(http.StatusUnauthorized, "error on extracting authentication token")
		}

		secretAuthKey := os.Getenv("SECRET_AUTH_KEY")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretAuthKey), nil
		})
		if err != nil {
			return c.JSON(http.StatusUnauthorized, "error on validating authentication token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.JSON(http.StatusUnauthorized, "error on validating authentication token")
		}

		c.Set("username", claims["username"])
		return handler(c)
	}
}
