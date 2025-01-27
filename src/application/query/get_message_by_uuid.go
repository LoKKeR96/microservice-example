package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/lokker96/microservice_example/domain/entity"
	"github.com/lokker96/microservice_example/domain/repository"
)

type GetMessageByUUIDQuery struct {
	messageRepository repository.MessageRepository
}

func NewGetMessageByUUIDQuery(
	messageRepository repository.MessageRepository,
) GetMessageByUUIDQuery {
	return GetMessageByUUIDQuery{
		messageRepository: messageRepository,
	}
}

func (g *GetMessageByUUIDQuery) Do(ctx context.Context, messageUUID uuid.UUID) (*entity.Message, error) {

	return g.messageRepository.GetMessageByUUID(messageUUID)
}
