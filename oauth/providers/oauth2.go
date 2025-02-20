package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"github.com/sjdaws/pkg/errors"
)

// Authenticator provider interface.
type Authenticator interface {
	Authenticate(ctx context.Context, code string, requestID string, verifier string) (*UserData, error)
	Authorise(verifier string) (*Request, error)
}

// OAuth2Provider interface.
type OAuth2Provider interface {
	GetUserdata() any
}

// OAuth2 root struct for all oauth2 providers.
type OAuth2 struct {
	Config  *oauth2.Config
	Profile Profile
}

// Profile information relating to retrieving user profile.
type Profile struct {
	Keys UserData
	URL  string
}

// Authenticate verifies a token and fetches user data.
func (o *OAuth2) Authenticate(ctx context.Context, code string, _ string, verifier string) (*UserData, error) {
	options := make([]oauth2.AuthCodeOption, 0)
	if verifier != "" {
		options = append(options, oauth2.VerifierOption(verifier))
	}

	token, err := o.Config.Exchange(ctx, code, options...)
	if err != nil {
		return nil, errors.Wrap(err, ErrInvalidToken)
	}

	return o.getUserData(ctx, o.Config.Client(ctx, token))
}

// Authorise creates an oauth request.
func (o *OAuth2) Authorise(verifier string) (*Request, error) {
	options := make([]oauth2.AuthCodeOption, 0)
	if verifier != "" {
		options = append(options, oauth2.S256ChallengeOption(verifier))
	}

	state := uuid.NewString()

	return &Request{
		ID:       uuid.NewString(),
		State:    state,
		URL:      o.Config.AuthCodeURL(state, options...),
		Verifier: verifier,
	}, nil
}

// getUserData uses the authentication token to retrieve user data.
func (o *OAuth2) getUserData(ctx context.Context, client *http.Client) (*UserData, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, o.Profile.URL, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrInvalidRequest)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, ErrInvalidUserData)
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode >= http.StatusBadRequest {
		return nil, errors.New(ErrInvalidResponseCode, response.StatusCode)
	}

	var data map[string]any

	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, errors.Wrap(err, ErrInvalidYAML)
	}

	return o.createUserData(data), nil
}

// createUserData parses the data from a token response into userdata based on provider mappings.
func (o *OAuth2) createUserData(data map[string]any) *UserData {
	var email string
	if o.Profile.Keys.Email != "" && data[o.Profile.Keys.Email] != nil {
		email = fmt.Sprintf("%v", data[o.Profile.Keys.Email])
	}

	var name string
	if o.Profile.Keys.Name != "" && data[o.Profile.Keys.Name] != nil {
		name = fmt.Sprintf("%v", data[o.Profile.Keys.Name])
	}

	var userID string
	if o.Profile.Keys.ID != "" && data[o.Profile.Keys.ID] != nil {
		userID = fmt.Sprintf("%v", data[o.Profile.Keys.ID])
	}

	var username string
	if o.Profile.Keys.Username != "" && data[o.Profile.Keys.Username] != nil {
		username = fmt.Sprintf("%v", data[o.Profile.Keys.Username])
	}

	return &UserData{
		Email:    email,
		ID:       userID,
		Name:     name,
		Username: username,
	}
}
