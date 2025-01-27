package controller

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/lokker96/microservice_example/application/query"
	"github.com/lokker96/microservice_example/infrastructure/container"
	"github.com/lokker96/microservice_example/infrastructure/controller/response"

	"github.com/labstack/echo/v4"
)

type GetMessagesByFilter struct{}

func NewGetMessagesByFilter() GetMessagesByFilter {
	return GetMessagesByFilter{}
}

type getMessagesByFilterRequest struct {
	Text string `query:"text" validate:"omitempty,max=255"`
	UUID string `query:"uuid" validate:"omitempty,max=255"`
}

func (sf *GetMessagesByFilter) Filter(ctx echo.Context, c container.Container) error {
	ctxBg := context.Background()

	var req getMessagesByFilterRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "error on binding filter parameters")
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "error on validating request")
	}

	queryRequest := query.GetMessagesByFilterRequest{}

	if req.Text != "" {
		queryRequest.Text = req.Text
	}

	if req.UUID != "" {
		parsedUUID, err := uuid.Parse(req.UUID)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, "error parsing uuid")
		}

		queryRequest.UUID = &parsedUUID
	}

	getAllMessagesQuery := c.GetMessagesByFilterQuery(ctxBg)

	messages, err := getAllMessagesQuery.Do(ctxBg, queryRequest)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "error on retrieving messages")
	}

	return ctx.JSON(http.StatusOK, response.NewMessageGroupResponse(messages))
}
