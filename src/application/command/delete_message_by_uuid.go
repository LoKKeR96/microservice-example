package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/lokker96/microservice_example/domain/repository"
)

type DeleteMessageCommand struct {
	ctx               context.Context
	messageRepository repository.MessageRepository
}

func NewDeleteMessageCommand(
	ctx context.Context,
	messageRepository repository.MessageRepository,
) DeleteMessageCommand {
	return DeleteMessageCommand{
		ctx:               ctx,
		messageRepository: messageRepository,
	}
}

func (ds *DeleteMessageCommand) Do(messageUUID uuid.UUID) error {
	return ds.messageRepository.DeleteMessageByUUID(messageUUID)
}
