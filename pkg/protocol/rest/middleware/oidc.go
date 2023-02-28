package middleware

import (
	"net/http"

	"github.com/muzcategui1106/kitchen-wizard/pkg/util"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/oidc"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/swagger"

	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// OauthCallBackHandler is used to handle oidc callback workflow
type OauthCallBackHandler struct {
	oauth2Config    oauth2.Config
	idTokenVerifier gooidc.IDTokenVerifier
}

func NewCallbackHandler(oauth2Config oauth2.Config, idTokenVerifier gooidc.IDTokenVerifier) *OauthCallBackHandler {
	return &OauthCallBackHandler{
		oauth2Config:    oauth2Config,
		idTokenVerifier: idTokenVerifier,
	}
}

// AddOIDCAuth adds oidc authentication workflow on all http handlers except healthz
func AddOIDCAuth(oauth2Config oauth2.Config, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// We do not want to be spammed by Kubernetes health check.
		// Do not log Kubernetes health check.
		// You can change this behavior as you wish.
		if r.Header.Get("X-Liveness-Probe") == "Healthz" {
			h.ServeHTTP(w, r)
			return
		}

		unauthenticatedPaths := []string{swagger.UIPrefix, oidc.CallbackURI}
		for _, path := range unauthenticatedPaths {
			if r.URL.Path == path {
				h.ServeHTTP(w, r)
				return
			}
		}

		state := util.RandStringRunes(16)
		http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
	})
}

// ServeHTTP handles the callback from an oidc auth flow
func (h *OauthCallBackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := LoggerFromContext(ctx)

	// Verify state.
	//state := r.URL.Query().Get("state")

	oauth2Token, err := h.oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		logger.Sugar().Errorf("error during oauth token exchange %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		logger.Sugar().Error("could not extract id token from request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Parse and verify ID Token payload.
	idToken, err := h.idTokenVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		logger.Sugar().Errorf("token could not be verified %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Extract custom claims.
	var claims struct {
		Email    string   `json:"email"`
		Verified bool     `json:"email_verified"`
		Groups   []string `json:"groups"`
	}
	if err := idToken.Claims(&claims); err != nil {
		logger.Sugar().Errorf("could  not extract claims from idToken %v", idToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	logger.Sugar().Infof("succesfully logged user with email %s and is verified %t", claims.Email, claims.Verified)
}
