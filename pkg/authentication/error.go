package authentication

// Error describes an authentication error against Keycloak.
type Error struct {
	err   error
	Text  string
	Token string
}

func (e Error) Error() string {
	return e.Text
}

// Unwrap returns the wrapped error.
func (e Error) Unwrap() error {
	return e.err
}

// NewError creates a new authentication error.
func NewError(err error, text string, token string) Error {
	return Error{err, text, token}
}