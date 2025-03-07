package error

import "fmt"

// For further development, it can be useful to look at:
// https://dev.to/tigorlazuardi/go-creating-custom-error-wrapper-and-do-proper-error-equality-check-11k7

type MessageAlreadyExists struct {
	UUID string
	error
}

func NewMessageAlreadyExists(uuid string) *MessageAlreadyExists {
	return &MessageAlreadyExists{
		UUID:  uuid,
		error: fmt.Errorf("message already exists: %s", uuid),
	}
}
