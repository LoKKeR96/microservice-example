package resolvers

import (
	"github.com/lokker96/microservice_example/domain/entity"
	"github.com/lokker96/microservice_example/infrastructure/graph/model"
)

func NewMessageResponse(message entity.Message) *model.Message {
	return &model.Message{
		Text: &message.Text,
		UUID: message.UUID.String(),
	}
}
