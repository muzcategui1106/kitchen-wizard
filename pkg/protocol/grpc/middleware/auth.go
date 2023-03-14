package middleware

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/oidc"
	"google.golang.org/grpc/metadata"
)

func WithOauth() runtime.ServeMuxOption {
	return runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
		md := metadata.Pairs()
		for _, headerKey := range []string{oidc.EmailKey, oidc.AccessTokenKey, oidc.RefreshTokenKey} {
			header := request.Header.Get(headerKey)
			md.Append(headerKey, header)
		}
		return md
	})
}
