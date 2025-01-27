package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/lokker96/microservice_example/domain/service"
	"github.com/palantir/stacktrace"

	"reflect"
)

type UpdateMessageByUUIDCommand struct {
	ctx                  context.Context
	messageEditorService service.MessageEditor
}

type UpdateMessageByUUIDRequest struct {
	Text *string
}

func NewUpdateMessageByUUIDCommand(
	ctx context.Context,
	messageEditorService service.MessageEditor) UpdateMessageByUUIDCommand {
	return UpdateMessageByUUIDCommand{
		ctx:                  ctx,
		messageEditorService: messageEditorService,
	}
}

func (us *UpdateMessageByUUIDCommand) Do(messageUUID uuid.UUID, messageRequest UpdateMessageByUUIDRequest) error {
	updates := make(map[string]interface{})

	// use reflection to get fields of the messageRequest struct
	requestValue := reflect.ValueOf(messageRequest)
	requestType := requestValue.Type()

	// iterate over each field in the struct
	for i := 0; i < requestValue.NumField(); i++ {
		// get the current field's value and its type
		currentFieldValue := requestValue.Field(i)
		currentFieldType := requestType.Field(i)

		if !currentFieldValue.IsNil() {
			updates[currentFieldType.Name] = currentFieldValue.Elem().Interface()
		}
	}

	if err := us.messageEditorService.Edit(messageUUID, updates); err != nil {
		return stacktrace.Propagate(err, "error on editing message by id")
	}

	return nil
}
