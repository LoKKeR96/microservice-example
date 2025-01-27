package controller

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/lokker96/microservice_example/application/command"
	"github.com/lokker96/microservice_example/infrastructure/container"
	"github.com/lokker96/microservice_example/infrastructure/controller/response"

	"github.com/labstack/echo/v4"
)

type UpdateMessage struct{}

func NewUpdateMessage() UpdateMessage {
	return UpdateMessage{}
}

type UpdateMessageRequest struct {
	Text *string `json:"text" validate:"omitempty,max=255"`
}

func (us *UpdateMessage) Update(ctx echo.Context, c container.Container) error {
	ctxBg := context.Background()
	var req UpdateMessageRequest

	parsedUUID, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "error parsing uuid")
	}

	if err = ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "error on binding message")
	}

	if err = ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "error on validating request")
	}

	updateMessageCommand := c.GetUpdateMessageByUUIDCommand(ctxBg)
	commandRequest := command.UpdateMessageByUUIDRequest{
		Text: req.Text,
	}
	if err = updateMessageCommand.Do(parsedUUID, commandRequest); err != nil {
		return ctx.JSON(http.StatusInternalServerError, "error on updating message")
	}

	return ctx.JSON(http.StatusOK, response.NewSuccessResponse())
}
