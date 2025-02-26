package plex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"

	"github.com/sjdaws/pkg/errors"
	"github.com/sjdaws/pkg/oauth/providers"
)

// Plex is the authentication provider for plex.
type Plex struct {
	CallbackURL  string
	ClientID     string
	ClientSecret string
	Endpoints    Endpoints
}

// Endpoints to communicate with plex api.
type Endpoints struct {
	AuthURL      string
	CreatePinURL string
	UserURL      string
	VerifyPinURL string
}

// New a new authenticator instance.
func New(callbackURL string, clientID string, clientSecret string, _ map[string]string) *Plex {
	return &Plex{
		CallbackURL:  callbackURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoints: Endpoints{
			AuthURL:      "https://app.plex.tv/auth",
			CreatePinURL: "https://plex.tv/api/v2/pins",
			UserURL:      "https://plex.tv/api/v2/user",
			VerifyPinURL: "https://plex.tv/api/v2/pins/%s",
		},
	}
}

// Authenticate verifies a token and fetches user data.
func (a *Plex) Authenticate(ctx context.Context, code string, requestID string, verifier string) (*providers.UserData, error) {
	var authToken struct {
		Token string `json:"authToken"`
	}

	data := bytes.NewBufferString(fmt.Sprintf(`code=%s&X-Plex-Client-Identifier=%s`, code, requestID))

	err := a.request(ctx, http.MethodGet, fmt.Sprintf(a.Endpoints.VerifyPinURL, verifier), data, &authToken)
	if err != nil {
		return nil, errors.Wrap(err, providers.ErrInvalidToken)
	}

	var user struct {
		ID       int    `json:"id"`
		Email    string `json:"email"`
		Name     string `json:"friendlyName"`
		Username string `json:"username"`
	}

	data = bytes.NewBufferString(
		fmt.Sprintf(
			`X-Plex-Client-Identifier=%s&X-Plex-Product=%s&X-Plex-Token=%s`,
			requestID,
			a.ClientSecret,
			authToken.Token,
		),
	)

	err = a.request(ctx, http.MethodGet, a.Endpoints.UserURL, data, &user)
	if err != nil {
		return nil, errors.Wrap(err, providers.ErrInvalidUserData)
	}

	return &providers.UserData{
		ID:       strconv.Itoa(user.ID),
		Email:    user.Email,
		Name:     user.Name,
		Username: user.Username,
	}, nil
}

// Authorise creates an oauth request.
func (a *Plex) Authorise(_ string) (*providers.Request, error) {
	var pin struct {
		Code string `json:"code"`
		ID   int    `json:"id"`
	}

	requestID := uuid.NewString()

	data := bytes.NewBufferString(
		fmt.Sprintf(
			`strong=true&X-Plex-Client-Identifier=%s&X-Plex-Product=%s`,
			requestID,
			a.ClientSecret,
		),
	)

	err := a.request(context.Background(), http.MethodPost, a.Endpoints.CreatePinURL, data, &pin)
	if err != nil {
		return nil, errors.Wrap(err, "pin creation failed")
	}

	uri, _ := url.Parse(a.Endpoints.AuthURL)

	query := uri.Query()
	query.Set("clientID", requestID)
	query.Set("code", pin.Code)
	query.Set("forwardUrl", a.CallbackURL)
	query.Set("context[device][product]", a.ClientSecret)

	return &providers.Request{
		ID:       requestID,
		State:    pin.Code,
		URL:      fmt.Sprintf("%s#?%s", uri.String(), query.Encode()),
		Verifier: strconv.Itoa(pin.ID),
	}, nil
}

// request perform a request to the plex API.
func (a *Plex) request(ctx context.Context, method string, endpoint string, body io.Reader, decode any) error {
	request, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return errors.Wrap(err, providers.ErrInvalidRequest)
	}

	request.Header.Add("Content-Type", "multipart/form-data")
	request.Header.Add("Accept", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "invalid response received")
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode >= http.StatusBadRequest {
		return errors.Wrap(err, providers.ErrInvalidResponseCode, response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(&decode); err != nil {
		return errors.Wrap(err, providers.ErrInvalidYAML)
	}

	return nil
}
