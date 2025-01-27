package repository

import (
	"github.com/google/uuid"
	"github.com/lokker96/microservice_example/domain/entity"
)

type MessageRepository interface {
	CreateMessage(message *entity.Message) error
	DeleteMessageByUUID(messageUUID uuid.UUID) error
	GetMessagesByFilter(messageText string, messageUUID *uuid.UUID) ([]*entity.Message, error)
	GetMessageByUUID(messageUUID uuid.UUID) (*entity.Message, error)
	UpdateMessage(message *entity.Message) error
	GetMessageByID(id uint) (*entity.Message, error)
	UpdateMessageFieldsByMessage(message *entity.Message, updates map[string]interface{}) error
}
