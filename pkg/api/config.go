package api

import "github.com/muzcategui1106/kitchen-wizard/pkg/util/oidc"

// Config is configuration for Server
type Config struct {
	OidcProviderConfig oidc.ProviderConfig
}
