package providers

// Request contains information to perform an oauth request.
type Request struct {
	ID       string
	State    string
	URL      string
	Verifier string
}

// UserData represents user data retrieved from a provider.
type UserData struct {
	Email    string
	ID       string
	Name     string
	Username string
}

const (
	// ErrInvalidRequest error when unable to create request.
	ErrInvalidRequest = "unable to create request"

	// ErrInvalidResponseCode error when bad response code is not ~200.
	ErrInvalidResponseCode = "invalid response code received: %d"

	// ErrInvalidToken error when token exchange fails.
	ErrInvalidToken = "token exchange failed"

	// ErrInvalidUserData error when retrieving user data fails.
	ErrInvalidUserData = "user data retrieval failed"

	// ErrInvalidYAML error when decoding yaml fails.
	ErrInvalidYAML = "yaml decode failed"
)
