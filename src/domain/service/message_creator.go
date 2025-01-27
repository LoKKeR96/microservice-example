package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/lokker96/microservice_example/domain/entity"
	domainErr "github.com/lokker96/microservice_example/domain/error"
	"github.com/lokker96/microservice_example/domain/repository"
	"github.com/palantir/stacktrace"
)

type MessageCreator interface {
	Create(message entity.Message) error
}

type messageCreator struct {
	ctx               context.Context
	messageRepository repository.MessageRepository
}

func NewMessageCreator(ctx context.Context, messageRepository repository.MessageRepository) MessageCreator {
	return &messageCreator{
		ctx:               ctx,
		messageRepository: messageRepository,
	}
}

func (sc *messageCreator) Create(message entity.Message) error {

	newUUID, err := uuid.NewUUID()
	if err != nil {
		return stacktrace.Propagate(err, "error creating message")
	}

	message.UUID = newUUID

	// if successful, early return
	err = sc.messageRepository.CreateMessage(&message)
	if err == nil {
		return nil
	}

	existingMessage, _ := sc.messageRepository.GetMessageByUUID(message.UUID)
	if existingMessage != nil {
		return stacktrace.Propagate(
			domainErr.NewMessageAlreadyExists(message.UUID.String()), //custom error handling in domain error
			"error on creating message: already exists",
		)
	}

	return stacktrace.Propagate(
		err,
		"error on creating message: %s",
		message.Text,
	)
}
