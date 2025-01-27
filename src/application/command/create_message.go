package command

import (
	"context"

	"github.com/lokker96/microservice_example/domain/entity"
	"github.com/lokker96/microservice_example/domain/service"
	"github.com/palantir/stacktrace"
)

type CreateMessageCommand struct {
	ctx                   context.Context
	messageCreatorService service.MessageCreator
}

type CreateMessageRequest struct {
	Text *string
}

func NewCreateMessageCommand(
	ctx context.Context,
	messageCreatorService service.MessageCreator) CreateMessageCommand {
	return CreateMessageCommand{
		ctx:                   ctx,
		messageCreatorService: messageCreatorService,
	}
}

func (cs *CreateMessageCommand) Do(messageRequest CreateMessageRequest) error {
	// constructing the entity from the bound request passed from the controller
	message := entity.Message{
		Text: *messageRequest.Text,
	}

	// calling domain message creator service
	if err := cs.messageCreatorService.Create(message); err != nil {
		return stacktrace.Propagate(err, "error on creating message")
	}

	return nil
}
