package restclient

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

const authTimeout = 60 * time.Second

type identityV2AuthenticatorImpl struct {
	username string
	password string
	apikey   string

	restClient *Client

	token           string
	tokenExpiration time.Time
}

// IdentityV2Authenticator provides an implementation of the Rackspace Cloud Identity v2.0
// authentication flow.
// The identityUrl should be the base URL of the Identity endpoint, such as "https://identity.api.rackspacecloud.com".
// Either password or apikey can be provided with the other passed an empty string.
//
// Info about Identity v2.0 is available at https://developer.rackspace.com/docs/cloud-identity/v2/
func IdentityV2Authenticator(identityUrl string, username string, password string, apikey string) (Interceptor, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	if password == "" && apikey == "" {
		return nil, errors.New("password or Apikey is required")
	}

	// looks slightly convoluted, but dogfood our own library to access the Identity REST API
	restClient := New()
	err := restClient.SetBaseUrl(identityUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid Identity URL: %w", err)
	}
	restClient.Timeout = authTimeout

	impl := &identityV2AuthenticatorImpl{
		username:   username,
		password:   password,
		apikey:     apikey,
		restClient: restClient,
	}

	return impl.intercept, nil
}

type identityAuthApikeyReq struct {
	Auth struct {
		Credentials struct {
			Username string `json:"username"`
			Apikey   string `json:"apiKey"`
		} `json:"RAX-KSKEY:apiKeyCredentials"`
	} `json:"auth"`
}

type identityAuthPasswordReq struct {
	Auth struct {
		Credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"passwordCredentials"`
	} `json:"auth"`
}

// identityAuthResp only picks out the fields needed and ignores the majority of response content
type identityAuthResp struct {
	Access struct {
		Token struct {
			Id      string
			Expires time.Time
		}
	}
}

func (a *identityV2AuthenticatorImpl) intercept(req *http.Request, next NextCallback) (*http.Response, error) {
	if time.Now().After(a.tokenExpiration) {
		if err := a.authenticate(); err != nil {
			return nil, err
		}
	}

	// inject the auth token into the user's REST request
	req.Header.Set("x-auth-token", a.token)

	return next(req)
}

func (a *identityV2AuthenticatorImpl) authenticate() error {

	var req interface{}
	if a.apikey != "" {
		auth := &identityAuthApikeyReq{}
		auth.Auth.Credentials.Username = a.username
		auth.Auth.Credentials.Apikey = a.apikey
		req = auth
	} else {
		auth := &identityAuthPasswordReq{}
		auth.Auth.Credentials.Username = a.username
		auth.Auth.Credentials.Password = a.password
		req = auth
	}

	var resp identityAuthResp

	err := a.restClient.Exchange("POST", "/v2.0/tokens", nil,
		NewJsonEntity(req), NewJsonEntity(&resp))
	if err != nil {
		return fmt.Errorf("failed to issue token request: %w", err)
	}

	a.token = resp.Access.Token.Id
	a.tokenExpiration = resp.Access.Token.Expires

	return nil
}
