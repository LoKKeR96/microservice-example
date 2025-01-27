package response

import "github.com/lokker96/microservice_example/domain/entity"

type MessageResponse struct {
	Text string `json:"text"`
	UUID string `json:"uuid"`
}

// NewMessageResponse makes the response logic on the controllers cleaner and also strips any sensitive data from the response
// although in a typical production environment, this would be handled in the application layer.
func NewMessageResponse(message entity.Message) MessageResponse {
	return MessageResponse{
		Text: message.Text,
		UUID: message.UUID.String(),
	}
}

// NewMessageGroupResponse returns an array of message, this is used specifically by the Filter controller
func NewMessageGroupResponse(messages []*entity.Message) map[string]interface{} {
	var responses []MessageResponse
	for _, message := range messages {
		// de-referencing and appending as repository returns pointers to handle nil returns
		responses = append(responses, NewMessageResponse(*message))
	}
	return map[string]interface{}{
		"data": responses,
	}
}
