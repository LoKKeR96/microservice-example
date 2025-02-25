package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lokker96/microservice_example/domain/entity"
	"github.com/lokker96/microservice_example/domain/repository"
	"github.com/palantir/stacktrace"
	"gorm.io/gorm"
)

type messageRepository struct {
	ctx context.Context // Intregation with middleware for handling timeouts and more
	db  *gorm.DB
}

func NewMessageRepository(ctx context.Context, db *gorm.DB) repository.MessageRepository {
	return &messageRepository{
		ctx: ctx,
		db:  db,
	}
}

func (r *messageRepository) CreateMessage(message *entity.Message) error {
	// using gorm transactions to make it easier to rollback if there are any issues, this is typically used in more complex repository methods to preserve data integrity
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(r.ctx).Create(&message).Error; err != nil {
			return stacktrace.Propagate(err, "error on creating message in db")
		}

		return nil
	})
}

func (r *messageRepository) DeleteMessageByUUID(messageUUID uuid.UUID) error {
	result := r.db.WithContext(r.ctx).Where("uuid = ?", messageUUID.String()).Delete(&entity.Message{})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return stacktrace.Propagate(
				fmt.Errorf("message not found with uuid: %s", messageUUID.String()),
				"error on deleting message by uuid",
			)
		}
		return stacktrace.Propagate(result.Error, "error on deleting message by uuid %s", messageUUID.String())
	}

	if result.RowsAffected == 0 {
		return stacktrace.Propagate(
			fmt.Errorf("message not found with uuid: %s", messageUUID.String()),
			"error on deleting message by uuid",
		)
	}

	return nil
}

func (r *messageRepository) GetMessagesByFilter(messageText string, messageUUID *uuid.UUID) ([]*entity.Message, error) {
	var messages []*entity.Message

	queryBuilder := r.db.WithContext(r.ctx)

	filters := map[string]string{
		"text": messageText,
	}
	if messageUUID != nil {
		filters["uuid"] = messageUUID.String()
	}

	for column, value := range filters {
		if value != "" {
			queryBuilder = queryBuilder.Where(column+" = ?", value)
		}
	}

	result := queryBuilder.Find(&messages)
	if result.Error != nil {
		return nil, stacktrace.Propagate(result.Error, "error on getting all messages")
	}

	return messages, nil
}

func (r *messageRepository) GetMessageByUUID(messageUUID uuid.UUID) (*entity.Message, error) {
	var message entity.Message

	result := r.db.WithContext(r.ctx).Where(&entity.Message{
		UUID: messageUUID,
	}).First(&message)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, stacktrace.Propagate(result.Error, "error on retrieving message by name")
	}

	return &message, nil
}

func (r *messageRepository) UpdateMessage(message *entity.Message) error {
	result := r.db.WithContext(r.ctx).Save(message)
	if result.Error != nil {
		return stacktrace.Propagate(result.Error, "error on updating message")
	}

	return nil
}

func (r *messageRepository) GetMessageByID(id uint) (*entity.Message, error) {
	var message entity.Message

	result := r.db.WithContext(r.ctx).First(&message, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, stacktrace.Propagate(result.Error, "error on retrieving message by id")
	}

	return &message, nil
}

func (r *messageRepository) UpdateMessageFieldsByMessage(message *entity.Message, updates map[string]interface{}) error {
	result := r.db.Model(message).Updates(updates)
	if result.Error != nil {
		return stacktrace.Propagate(result.Error, "error on updating message fields")
	}

	return nil
}
