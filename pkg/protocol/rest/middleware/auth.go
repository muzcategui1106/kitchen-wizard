package middleware

import (
	"context"
	"encoding/gob"
	"net/http"

	"github.com/muzcategui1106/kitchen-wizard/pkg/util"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/oidc"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/swagger"

	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

func init() {
	gob.Register(&gooidc.IDToken{})
}

// constants for URL paths
const (
	BaseAuthPathV1 = "/auth/v1/"
)

// cookie and sesssion constants
const (
	IDTokenKey     = "id_token"
	UserSessionKey = "user_session"
)

// variables for URL paths
var (
	v1Login = BaseAuthPathV1 + "login"
)

// AuthHandler is used to handle oidc callback workflow
// the handler provides a session store to ensure client browsers store and send cookies containing the access/id/refresh tokens
// TODO
// at this time the auth handler will only keep a session for the current server. Note that a more distributed approach perhaps store the session
// in a redis cache or in a database for later retrieval and verification
type AuthHandler struct {
	oauth2Config    oauth2.Config
	idTokenVerifier gooidc.IDTokenVerifier
	sessionStore    sessions.Store
}

// NewAuthHandler handles all authentication calls
func NewAuthHandler(oauth2Config oauth2.Config, idTokenVerifier gooidc.IDTokenVerifier, sessionKey []byte) *AuthHandler {
	sessionStore := sessions.NewCookieStore(sessionKey)
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
	}

	return &AuthHandler{
		oauth2Config:    oauth2Config,
		idTokenVerifier: idTokenVerifier,
		sessionStore:    sessionStore,
	}
}

// AuthenticationInterceptor adds oidc authentication workflow on all http handlers except healthz
func (auh *AuthHandler) AuthenticationInterceptor(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := LoggerFromContext(ctx)

		// We do not want to be spammed by Kubernetes health check.
		// Do not log Kubernetes health check.
		// You can change this behavior as you wish.
		if r.Header.Get("X-Liveness-Probe") == "Healthz" {
			h.ServeHTTP(w, r)
			return
		}

		unauthenticatedPaths := []string{swagger.UIPrefix, oidc.CallbackURI, v1Login, "/api/v1/healthz"}
		for _, path := range unauthenticatedPaths {
			if r.URL.Path == path {
				h.ServeHTTP(w, r)
				return
			}
		}

		// do not do login if a session ID has been extracted
		session, err := auh.sessionStore.Get(r, UserSessionKey)
		if err != nil {
			logger.Sugar().Errorf("could not get or create session due to %s", err)
			goto doLogin
		}

		if session.IsNew {
			logger.Info("new sessions detected, redirecting to login")
			goto doLogin
		} else {
			idToken, ok := session.Values[IDTokenKey]
			if !ok {
				logger.Sugar().Error("id token not found in session this should not happen, doing login")
				goto doLogin
			} else {
				h.ServeHTTP(w, r.WithContext(context.WithValue(ctx, IDTokenKey, idToken)))
				return
			}
		}

	doLogin:
		logger.Sugar().Debug("user does nnot have an existing session redirecting to login")
		http.Redirect(w, r, v1Login, http.StatusFound)
	})
}

// ServeHTTP handles the callback from an oidc auth flow
func (auh *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case oidc.CallbackURI:
		auh.handleOIDCCallback(w, r)
	case v1Login:
		auh.handleV1Login(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (auh *AuthHandler) handleV1Login(w http.ResponseWriter, r *http.Request) {
	state := util.RandStringRunes(16)
	nonce := util.RandStringRunes(16)

	http.Redirect(w, r, auh.oauth2Config.AuthCodeURL(state, gooidc.Nonce(nonce)), http.StatusFound)
}

func (auh *AuthHandler) handleOIDCCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := LoggerFromContext(ctx)

	oauth2Token, err := auh.oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		logger.Sugar().Errorf("error during oauth token exchange %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra(IDTokenKey).(string)
	if !ok {
		logger.Sugar().Error("could not extract id token from request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Parse and verify ID Token payload.
	idToken, err := auh.idTokenVerifier.Verify(ctx, rawIDToken)
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

	// create a session for the user
	session, err := auh.sessionStore.Get(r, UserSessionKey)
	if err != nil {
		logger.Sugar().Errorf("could not get session from session store. error was .... %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// store the id token in the session
	session.Values[IDTokenKey] = idToken
	err = session.Save(r, w)
	if err != nil {
		logger.Sugar().Error("unable to save session due to ... %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Sugar().Infof("succesfully logged user with email %s and is verified %t", claims.Email, claims.Verified)
	w.Write([]byte("user has logged in"))
}
