package errors

// Response is the default JSON error response structure returned if somethins went wrong.
type Response struct {
	Error string `json:"error"`
}

// NewResponse creates a new error response.
func NewResponse(message string) Response {
	return Response{message}
}