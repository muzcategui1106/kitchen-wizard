package oidc

import (
	"context"

	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/muzcategui1106/kitchen-wizard/pkg/logger"
	"golang.org/x/oauth2"
)

const (
	CallbackURI = "/v1/auth/oidc/callback"
)

// ProviderConfig contains all the necessary configuration OIDC
type ProviderConfig struct {
	// ProviderURL for the OIDC provider to be used
	ProviderURL string

	// OidcClientID it is the client ID to use for oidc authentication
	OidcClientID string

	// OidcClientSecret is the secret to be used to communicate to dex for any oauth workflow
	OidcClientSecret string

	// OidcRedirectURL is the URL to be used for oidc callback. The path will be harcoded as a constant in this codebase. Perhaps
	// it is better to just include everything here but at this time it makes more sense to me to separate them.
	OidcRedirectURL string
}

// CreateOIDCClient creates a proper OIDC client
func CreateOIDCClient(ctx context.Context, conf ProviderConfig) (oauth2.Config, *gooidc.IDTokenVerifier, error) {
	logger := logger.Log

	// Initialize a provider by specifying dex's issuer URL.
	provider, err := gooidc.NewProvider(ctx, conf.ProviderURL)
	if err != nil {
		logger.Sugar().Errorf("could not create oidc provider", err)
		return oauth2.Config{}, nil, err
	}

	// Configure the OAuth2 config with the client values.
	oauth2Config := oauth2.Config{
		// client_id and client_secret of the client.
		ClientID:     conf.OidcClientID,
		ClientSecret: conf.OidcClientSecret,

		// The redirectURL.
		RedirectURL: conf.OidcRedirectURL + CallbackURI,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		//
		// Other scopes, such as "groups" can be requested.
		Scopes: []string{gooidc.ScopeOpenID, "profile", "email", "groups"},
	}

	// Create an ID token parser.
	return oauth2Config, provider.Verifier(&gooidc.Config{ClientID: "example-app"}), nil
}
