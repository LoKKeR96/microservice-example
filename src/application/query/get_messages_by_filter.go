package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/lokker96/microservice_example/domain/entity"
	"github.com/lokker96/microservice_example/domain/repository"
	"github.com/palantir/stacktrace"
)

type GetMessagesByFilterQuery struct {
	messageRepository repository.MessageRepository
}

type GetMessagesByFilterRequest struct {
	Text string
	UUID *uuid.UUID
}

func NewGetMessagesByFilterQuery(
	messageRepository repository.MessageRepository,
) GetMessagesByFilterQuery {
	return GetMessagesByFilterQuery{
		messageRepository: messageRepository,
	}
}

func (ga *GetMessagesByFilterQuery) Do(ctx context.Context, filter GetMessagesByFilterRequest) ([]*entity.Message, error) {
	allMessages, err := ga.messageRepository.GetMessagesByFilter(
		filter.Text,
		filter.UUID,
	)
	if err != nil {
		return nil, stacktrace.Propagate(err, "error on filtering messages")
	}
	return allMessages, nil
}
