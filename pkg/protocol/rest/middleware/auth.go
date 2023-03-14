package middleware

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	UserSessionKey = "session_id"
	IDTokenKey     = "id_token"
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
			email, ok := session.Values[oidc.EmailKey]
			if !ok {
				logger.Sugar().Error("session")
				goto doLogin
			} else {
				r.Header.Add(oidc.EmailKey, string(fmt.Sprintf("%v", email)))
				h.ServeHTTP(w, r)
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

// TODO this function saves all user info and refresh token in session as cookies. This is not correct
// we only need to pass the access token and we can store the rest in redis or a database or something along
// those lines. However we do not want to complicate ourselves with this at the moment
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

		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
	}
	if err := idToken.Claims(&claims); err != nil {
		logger.Sugar().Errorf("could  not extract claims from idToken %v", idToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Println(oauth2Token.AccessToken)
	fmt.Println(oauth2Token.Expiry)

	// proceed to get the user's first and last name
	req, err := http.NewRequest("GET", "https://api.linkedin.com/v2/userinfo", nil)
	if err != nil {
		logger.Sugar().Errorf("Error creating request to linkedin %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := auh.oauth2Config.Client(ctx, oauth2Token).Do(req)
	if err != nil {
		logger.Sugar().Errorf("Error sending request to linkedin %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// Read the response body and parse the JSON data into a User struct
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Sugar().Errorf("Error reading user response from linkedin %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	type FirstName struct {
		Localized map[string]string `json:"localized"`
	}

	type LastName struct {
		Localized map[string]string `json:"localized"`
	}

	type User struct {
		FirstName FirstName `json:"firstName"`
		LastName  LastName  `json:"lastName"`
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		logger.Sugar().Errorf("Error parsin json response from linkedin %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Print the user's first name and last name
	fmt.Println("First Name:", user.FirstName.Localized["en_US"])
	fmt.Println("Last Name:", user.LastName.Localized["en_US"])

	// create a session for the user
	session, err := auh.sessionStore.Get(r, UserSessionKey)
	if err != nil {
		logger.Sugar().Errorf("could not get session from session store. error was .... %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// store the id token in the session
	session.Values[oidc.EmailKey] = claims.Email
	session.Values[oidc.AccessTokenKey] = oauth2Token.AccessToken
	err = session.Save(r, w)
	if err != nil {
		logger.Sugar().Error("unable to save session due to ... %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Sugar().Infof("succesfully logged user with email %s and is verified %t", claims.Email, claims.Verified)
	w.Write([]byte("user has logged in"))
}
