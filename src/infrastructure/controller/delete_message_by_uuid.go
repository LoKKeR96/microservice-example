package controller

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/lokker96/microservice_example/infrastructure/container"
	"github.com/lokker96/microservice_example/infrastructure/controller/response"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type DeleteMessage struct{}

func NewDeleteMessage() DeleteMessage {
	return DeleteMessage{}
}

func (ds *DeleteMessage) Delete(ctx echo.Context, c container.Container) error {
	ctxBg := context.Background()

	parsedUUID, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "error parsing uuid")
	}

	deleteMessageCommand := c.GetDeleteMessageCommand(ctxBg)
	err = deleteMessageCommand.Do(parsedUUID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "not found") {
			return ctx.JSON(http.StatusNotFound, "message not found")
		}

		return ctx.JSON(http.StatusInternalServerError, "error deleting message")
	}

	return ctx.JSON(http.StatusOK, response.NewSuccessResponse())
}
