package route

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/lokker96/microservice_example/infrastructure/container"
	"github.com/lokker96/microservice_example/infrastructure/controller"

	middleware "github.com/lokker96/microservice_example/infrastructure/http"
)

func Routes(c container.Container) *echo.Echo {
	e := echo.New()

	// adding custom validator for modularity and scalability
	// separating this logic makes it easier to make validation changes in the future
	e.Validator = &CustomValidator{validator.New()}

	createMessageController := controller.NewCreateMessage()
	deleteMessageController := controller.NewDeleteMessage()
	getAllMessages := controller.NewGetMessagesByFilter()
	getMessageController := controller.NewGetMessage()
	updateMessageController := controller.NewUpdateMessage()

	userAuthenticationController := controller.NewUserAuthentication()

	// route definition

	e.POST(
		"/message/create",
		func(ctx echo.Context) error {
			return createMessageController.Create(ctx, c)
		},
		middleware.AuthenticationMiddleware,
	)

	e.PUT(
		"/message/:uuid/update",
		func(ctx echo.Context) error {
			return updateMessageController.Update(ctx, c)
		},
		middleware.AuthenticationMiddleware,
	)

	e.DELETE(
		"/message/:uuid/delete",
		func(ctx echo.Context) error {
			return deleteMessageController.Delete(ctx, c)
		},
		middleware.AuthenticationMiddleware,
	)

	e.POST(
		"/message/search",
		func(ctx echo.Context) error {
			return getAllMessages.Filter(ctx, c)
		},
		middleware.AuthenticationMiddleware,
	)

	e.GET(
		"/message/:uuid/get",
		func(ctx echo.Context) error {
			return getMessageController.Get(ctx, c)
		},
		middleware.AuthenticationMiddleware,
	)

	e.POST("/user/login",
		func(ctx echo.Context) error {
			return userAuthenticationController.Authenticate(ctx, c)
		},
	)

	return e
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
