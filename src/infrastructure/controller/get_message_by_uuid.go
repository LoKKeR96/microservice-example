package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/lokker96/microservice_example/infrastructure/container"
	"github.com/lokker96/microservice_example/infrastructure/controller/response"

	"github.com/labstack/echo/v4"

	"net/http"
)

type GetMessage struct{}

func NewGetMessage() GetMessage {
	return GetMessage{}
}

func (gs *GetMessage) Get(ctx echo.Context, c container.Container) error {
	ctxBg := context.Background()

	messageUUID := ctx.Param("uuid")

	parsedUUID, err := uuid.Parse(messageUUID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "error parsing uuid")
	}

	getMessageByUUIDQuery := c.GetMessageByUUIDQuery(ctxBg)

	message, err := getMessageByUUIDQuery.Do(ctxBg, parsedUUID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "error on retrieving message")
	}

	if message == nil {
		return ctx.JSON(http.StatusNotFound, "message not found")
	}

	return ctx.JSON(http.StatusOK, response.NewMessageResponse(*message))
}
