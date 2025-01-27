package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lokker96/microservice_example/domain/repository"
	"github.com/palantir/stacktrace"
)

type MessageEditor interface {
	Edit(messageUUID uuid.UUID, updates map[string]interface{}) error
}

type messageEditor struct {
	ctx               context.Context
	messageRepository repository.MessageRepository
}

func NewMessageEditor(ctx context.Context, messageRepository repository.MessageRepository) MessageEditor {
	return &messageEditor{
		ctx:               ctx,
		messageRepository: messageRepository,
	}
}

func (se *messageEditor) Edit(messageUUID uuid.UUID, updates map[string]interface{}) error {
	existingMessage, err := se.messageRepository.GetMessageByUUID(messageUUID)
	if err != nil {
		return stacktrace.Propagate(err, "error on getting message by uuid")
	}

	if existingMessage == nil {
		return stacktrace.Propagate(
			fmt.Errorf("message not found"),
			"error on editing message: not found",
		)
	}

	if len(updates) == 0 {
		return nil
	}

	err = se.messageRepository.UpdateMessageFieldsByMessage(existingMessage, updates)
	if err != nil {
		return stacktrace.Propagate(err, "error on updating message fields")
	}

	return nil
}
