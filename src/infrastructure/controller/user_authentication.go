package controller

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/lokker96/microservice_example/infrastructure/container"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/palantir/stacktrace"
)

const (
	USERNAME        = "member1"
	PASSWORD        = "password123"
	TokenExpiration = 17*time.Hour + 30*time.Minute
)

type UserAuthentication struct {
	SecretAuthKey string
}

func NewUserAuthentication() *UserAuthentication {
	return &UserAuthentication{
		SecretAuthKey: os.Getenv("SECRET_AUTH_KEY"),
	}
}

type authenticationCredentials struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Authenticate unused container to maintain consistency and maintainability for future development
func (ia *UserAuthentication) Authenticate(ctx echo.Context, c container.Container) error {
	var credentials authenticationCredentials

	if err := ctx.Bind(&credentials); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Error binding credentials")
	}

	if err := ctx.Validate(&credentials); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Error validating credentials")
	}

	token, err := ia.authenticate(credentials.Username, credentials.Password)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

// just a simple helper function to generate a jwt token
// in a production environment you'd have the authentication logic housed in a service on the domain layer
// which is called through an application/command
// for middleware demonstration purposes I've implemented the logic directly in the controller
func (ia *UserAuthentication) authenticate(username, password string) (string, error) {
	if username != USERNAME || password != PASSWORD {
		return "", stacktrace.Propagate(
			fmt.Errorf("error on validating credentials"),
			"error on validating credentials",
		)
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(TokenExpiration).Unix(),
	})

	return claims.SignedString([]byte(ia.SecretAuthKey))
}
