package authentication

// authenticationError describes an authentication error against Keycloak.
type authenticationError struct {
	err   error
	Text  string
	Token string
}

func (e authenticationError) Error() string {
	return e.Text
}

// Unwrap returns the wrapped error.
func (e authenticationError) Unwrap() error {
	return e.err
}

// newAuthenticationError creates a new authentication error.
func newAuthenticationError(text string, token string, err error) authenticationError {
	return authenticationError{err, text, token}
}