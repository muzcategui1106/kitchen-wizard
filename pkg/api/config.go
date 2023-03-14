package api

import (
	"database/sql"

	"github.com/muzcategui1106/kitchen-wizard/pkg/util/oidc"
)

// Config is configuration for Server
type Config struct {
	OidcProviderConfig oidc.ProviderConfig
	DBConn             *sql.DB
}
