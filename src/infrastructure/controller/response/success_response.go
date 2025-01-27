package response

type SuccessResponse struct {
	Success bool `json:"success"`
}

// NewSuccessResponse creating a method so that the response logic in the controller is cleaner
func NewSuccessResponse() SuccessResponse {
	return SuccessResponse{
		Success: true,
	}
}
