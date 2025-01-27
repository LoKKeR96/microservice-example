package error

import "fmt"

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
