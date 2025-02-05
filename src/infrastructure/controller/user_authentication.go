package controller

import (
	"net/http"

	"github.com/lokker96/microservice_example/infrastructure/container"

	"github.com/labstack/echo/v4"
)

type UserAuthentication struct{}

func NewUserAuthentication() *UserAuthentication {
	return &UserAuthentication{}
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

	userAuthenticationService := c.GetUserAuthenticationService()

	token, err := userAuthenticationService.Authenticate(credentials.Username, credentials.Password)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
