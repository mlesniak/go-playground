package errors

// Response is the default JSON error response structure returned if somethins went wrong.
type Response struct {
	Error string `json:"error"`
}
