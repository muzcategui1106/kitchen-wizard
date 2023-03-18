package middleware

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/oidc"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/swagger"

	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"golang.org/x/oauth2"
)

// constants for api group
const (
	authBasePath = "/auth"
)

// constants for versions
const (
	versionV1 = "/v1"
)

// constants for api paths
const (
	login = "/login"
)

func init() {
	gob.Register(&gooidc.IDToken{})
}

// cookie and sesssion constants
const (
	UserSessionKey = "session_id"
	IDTokenKey     = "id_token"
)

type V1LogiResponse struct {
	OK bool `json:"ok"`
}

// AuthHandler is used to handle oidc callback workflow
// the handler provides a session store to ensure client browsers store and send cookies containing the access/id/refresh tokens
// TODO
// at this time the auth handler will only keep a session for the current server. Note that a more distributed approach perhaps store the session
// in a redis cache or in a database for later retrieval and verification
type AuthHandler struct {
	oauth2Config    oauth2.Config
	idTokenVerifier gooidc.IDTokenVerifier
}

// NewAuthHandler handles all authentication calls
func NewAuthHandler(oauth2Config oauth2.Config, idTokenVerifier gooidc.IDTokenVerifier) *AuthHandler {
	return &AuthHandler{
		oauth2Config:    oauth2Config,
		idTokenVerifier: idTokenVerifier,
	}
}

func (auh *AuthHandler) AddAuthHandling(r *gin.Engine) {
	authGroup := r.Group(authBasePath + versionV1)
	authGroup.Any(login, auh.handleV1Login())
	authGroup.POST(oidc.CallbackURI, auh.handleOIDCCallback())
}

// AuthenticationInterceptor adds oidc authentication workflow on all http handlers except healthz
func (auh *AuthHandler) AuthenticationInterceptor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := LoggerFromContext(ctx)

		// We do not want to be spammed by Kubernetes health check.
		// Do not log Kubernetes health check.
		// You can change this behavior as you wish.
		if ctx.Request.Header.Get("X-Liveness-Probe") == "Healthz" {
			ctx.Next()
			return
		}

		unauthenticatedPaths := []string{swagger.UIPrefix, oidc.CallbackURI, authBasePath + versionV1 + login, "/api/v1/healthz"}
		for _, path := range unauthenticatedPaths {
			if ctx.Request.URL.Path == path {
				ctx.Next()
				return
			}
		}

		// do not do login if a session ID has been extracted
		session := sessions.DefaultMany(ctx, UserSessionKey)
		email := session.Get(oidc.EmailKey)
		if email == nil {
			logger.Sugar().Error("could not get session from request. redirecting to login")
			goto doLogin
		} else {
			ctx.Request.Header.Add(oidc.EmailKey, string(fmt.Sprintf("%v", email)))
			ctx.Next()
			return
		}

	doLogin:
		logger.Sugar().Debug("user does nnot have an existing session redirecting to login")
		ctx.Redirect(http.StatusFound, authBasePath+versionV1+login)
	}
}

func (auh *AuthHandler) handleV1Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		state := util.RandStringRunes(16)
		nonce := util.RandStringRunes(16)
		ctx.Redirect(http.StatusFound, auh.oauth2Config.AuthCodeURL(state, gooidc.Nonce(nonce)))
	}
}

// TODO this function saves all user info and refresh token in session as cookies. This is not correct
// we only need to pass the access token and we can store the rest in redis or a database or something along
// those lines. However we do not want to complicate ourselves with this at the moment
func (auh *AuthHandler) handleOIDCCallback() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := LoggerFromContext(ctx)

		codeAny, exists := ctx.Get("code")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "oidc did not return a code query parameter"})
			return
		}
		code := codeAny.(string)

		oauth2Token, err := auh.oauth2Config.Exchange(ctx, string(code))
		if err != nil {
			logger.Sugar().Errorf("error during oauth token exchange %s", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error during ouath token exchange"})
			return
		}

		// Extract the ID Token from OAuth2 token.
		rawIDToken, ok := oauth2Token.Extra(IDTokenKey).(string)
		if !ok {
			logger.Sugar().Error("could not extract id token from request")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not extract id token from request"})
			return
		}

		// Parse and verify ID Token payload.
		idToken, err := auh.idTokenVerifier.Verify(ctx, rawIDToken)
		if err != nil {
			logger.Sugar().Errorf("token could not be verified %s", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "oauth token could not be verified"})
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
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "could not extract claims from token"})
			return
		}

		fmt.Println(oauth2Token.AccessToken)
		fmt.Println(oauth2Token.Expiry)

		// proceed to get the user's first and last name
		req, err := http.NewRequest("GET", "https://api.linkedin.com/v2/userinfo", nil)
		if err != nil {
			logger.Sugar().Errorf("Error creating request to linkedin %s", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating request to linkedin"})
			return
		}

		res, err := auh.oauth2Config.Client(ctx, oauth2Token).Do(req)
		if err != nil {
			logger.Sugar().Errorf("Error sending request to linkedin %s", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending request to linkedin"})
			return
		}
		defer res.Body.Close()

		// Read the response body and parse the JSON data into a User struct
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Sugar().Errorf("Error reading user response from linkedin %s", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading user response from linkedin"})
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
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsin json response from linkedin"})
			return
		}

		// Print the user's first name and last name
		fmt.Println("First Name:", user.FirstName.Localized["en_US"])
		fmt.Println("Last Name:", user.LastName.Localized["en_US"])

		// do not do login if a session ID has been extracted
		session := sessions.DefaultMany(ctx, UserSessionKey)

		// store the id token in the session
		session.Set(oidc.EmailKey, claims.Email)
		session.Set(oidc.AccessTokenKey, oauth2Token.AccessToken)
		err = session.Save()
		if err != nil {
			logger.Sugar().Error("unable to save session due to ... %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unable to save session"})
			return
		}

		logger.Sugar().Infof("succesfully logged user with email %s and is verified %t", claims.Email, claims.Verified)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unable to save session"})
	}
}
