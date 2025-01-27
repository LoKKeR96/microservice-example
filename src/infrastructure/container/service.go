package container

import (
	"context"

	"github.com/lokker96/microservice_example/application/command"
	"github.com/lokker96/microservice_example/application/query"
	"github.com/lokker96/microservice_example/domain/repository"
	"github.com/lokker96/microservice_example/domain/service"
	"github.com/lokker96/microservice_example/infrastructure/persistence/mysql"
)

var _ Services = &Container{}

type Services interface {
	GetMessageRepository(ctx context.Context) repository.MessageRepository

	GetMessageCreatorService(ctx context.Context) service.MessageCreator
	GetMessageEditorService(ctx context.Context) service.MessageEditor

	GetCreateMessageCommand(ctx context.Context) command.CreateMessageCommand
	GetDeleteMessageCommand(ctx context.Context) command.DeleteMessageCommand
	GetUpdateMessageByUUIDCommand(ctx context.Context) command.UpdateMessageByUUIDCommand

	GetMessagesByFilterQuery(ctx context.Context) query.GetMessagesByFilterQuery
	GetMessageByUUIDQuery(ctx context.Context) query.GetMessageByUUIDQuery
}

// Repository

func (c *Container) GetMessageRepository(ctx context.Context) repository.MessageRepository {
	return mysql.NewMessageRepository(ctx, c.db)
}

// Service

func (c *Container) GetMessageCreatorService(ctx context.Context) service.MessageCreator {
	return service.NewMessageCreator(ctx, c.GetMessageRepository(ctx))
}

func (c *Container) GetMessageEditorService(ctx context.Context) service.MessageEditor {
	return service.NewMessageEditor(ctx, c.GetMessageRepository(ctx))
}

// Command

func (c *Container) GetCreateMessageCommand(ctx context.Context) command.CreateMessageCommand {
	return command.NewCreateMessageCommand(ctx, c.GetMessageCreatorService(ctx))
}

func (c *Container) GetDeleteMessageCommand(ctx context.Context) command.DeleteMessageCommand {
	return command.NewDeleteMessageCommand(ctx, c.GetMessageRepository(ctx))
}

func (c *Container) GetUpdateMessageByUUIDCommand(ctx context.Context) command.UpdateMessageByUUIDCommand {
	return command.NewUpdateMessageByUUIDCommand(ctx, c.GetMessageEditorService(ctx))
}

// Query

func (c *Container) GetMessagesByFilterQuery(ctx context.Context) query.GetMessagesByFilterQuery {
	return query.NewGetMessagesByFilterQuery(c.GetMessageRepository(ctx))
}

func (c *Container) GetMessageByUUIDQuery(ctx context.Context) query.GetMessageByUUIDQuery {
	return query.NewGetMessageByUUIDQuery(c.GetMessageRepository(ctx))
}
