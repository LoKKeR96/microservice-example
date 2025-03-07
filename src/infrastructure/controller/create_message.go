package controller

import (
	"context"
	"errors"

	"github.com/lokker96/microservice_example/application/command"
	domainErr "github.com/lokker96/microservice_example/domain/error"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lokker96/microservice_example/infrastructure/container"
	"github.com/lokker96/microservice_example/infrastructure/controller/response"
)

type CreateMessage struct{}

func NewCreateMessage() CreateMessage {
	return CreateMessage{}
}

type CreateMessageRequest struct {
	Text string `json:"text" validate:"required,max=255"`
}

// Create handles the creation of a message
func (cs *CreateMessage) Create(ctx echo.Context, c container.Container) error {
	ctxBg := context.Background()
	var req CreateMessageRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "error on binding message")
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "error on validating request")
	}

	// retrieve command from container through dependency injection
	createMessageCommand := c.GetCreateMessageCommand(ctxBg)
	commandRequest := command.CreateMessageRequest{
		Text: &req.Text,
	}

	if err := createMessageCommand.Do(commandRequest); err != nil {
		// handle custom domain error
		if errors.Is(err, &domainErr.MessageAlreadyExists{}) {
			// Ideally we should use a logger for errors and communicate
			// something more appropriate externally.
			return ctx.JSON(http.StatusConflict, err.Error())
		}
		return ctx.JSON(http.StatusInternalServerError, "error on creating message")
	}

	return ctx.JSON(http.StatusOK, response.NewSuccessResponse())
}
