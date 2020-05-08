package response

// Error is the default JSON error response structure returned if somethins went wrong.
type Error struct {
	Error string `json:"error"`
}

// NewError creates a new error response.
func NewError(message string) Error {
	return Error{message}
}