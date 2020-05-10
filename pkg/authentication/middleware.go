package authentication

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/mlesniak/go-demo/pkg/context"
	"github.com/mlesniak/go-demo/pkg/response"
	"github.com/rs/zerolog/log"
)

// KeycloakConfig defines configuration options for the middleware.
type KeycloakConfig struct {
	Protocol string
	Hostname string
	Port     string
	Realm    string

	IgnoredURL []string
	LoginURL   string
	LogoutURL  string
	RefreshURL string
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type logoutRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	RefreshToken string `json:"refreshToken"`
}

type authenticationResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// cache stores all already validated (and still valid) tokens.
// Before a token is authenticated, its expiration date is checked.
var cache map[string]int64 = make(map[string]int64)

// KeycloakWithConfig registers authentication endpoints and handles token validation on each request.
func KeycloakWithConfig(e *echo.Echo, config KeycloakConfig) func(next echo.HandlerFunc) echo.HandlerFunc {
	if config.Protocol == "" || config.Hostname == "" || config.Port == "" || config.Realm == "" || config.LoginURL == "" || config.LogoutURL == "" || config.RefreshURL == "" {
		panic("The keycloak configuration is invalid, at least one required property is empty.")
	}
	config.addEndpoints(e)

	// Check token on each request.
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(cc echo.Context) error {
			c := cc.(*context.CustomContext)
			log := c.Log()

			// Check for ignored URLS such as login routes.
			shouldURLbeChecked := true
			for _, i := range config.IgnoredURL {
				if i == c.Request().URL.RequestURI() {
					shouldURLbeChecked = false
					break
				}
			}
			if shouldURLbeChecked {
				// Has a valid token been submitted?
				if err := config.isAuthenticated(c); err != nil {
					e := err.(authenticationError)
					log.Info().Str("token", e.Token).Msg(e.Text)
					return c.JSON(http.StatusUnauthorized, response.NewError("Token is invalid"))
				}
				// Everything is ok, add authentication info to request context.
				addUserInfoToContext(c)
				// Update default context with username
				log = c.Log()
				log.Info().Msg("Successful authentication by token")
			}

			// Continue chain.
			if err := next(c); err != nil {
				c.Error(err)
			}
			return nil
		}
	}
}

// getKeycloakURLFor returns the fully qualified URL for the given operations based on the pre-defined configuration.
func (config *KeycloakConfig) getKeycloakURLFor(operation string) string {
	return fmt.Sprintf(
		"%s://%s:%s/auth/realms/%s/protocol/openid-connect/%s",
		config.Protocol, config.Hostname, config.Port, config.Realm, operation)
}

// IsAuthenticated returns true if the user submitted a valid JWT token.
func (config *KeycloakConfig) isAuthenticated(c *context.CustomContext) error {
	log := c.Log()
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return newAuthenticationError("Empty token", token, nil)
	}

	// Check cache before accessing keycloak.
	var jwtToken string
	if len(token) > 7 {
		jwtToken = token[7:]
	}
	expiresAt, found := cache[jwtToken]
	if found {
		now := time.Now().Unix()
		if expiresAt > now {
			log.Info().
				Str("token", jwtToken).
				Int64("expiresAt", expiresAt).
				Msg("Accepting token from cache")
			return nil
		}

		// Entry should be removed, since it's expired.
		log.Info().
			Str("token", token).
			Int64("expiresAt", expiresAt).
			Msg("Remove expired token from cache")
		delete(cache, token)
	}

	// Access userinfo in keycloak using the provided token to check if the token is valid.
	req, err := http.NewRequest("GET", config.getKeycloakURLFor("userinfo"), nil)
	req.Header.Add("Authorization", token)
	cl := &http.Client{}
	resp, err := cl.Do(req)

	if err != nil {
		return newAuthenticationError("Unable to get userinfo", token, err)
	}
	if resp.StatusCode == 200 {
		// Happy flow.
		// Add token on successful authentication w/o cache here, too. This is necessary
		// if the server has been restarted in the meantime while the user has still an active token.
		addTokenToCache(jwtToken)
		return nil
	}
	if resp.StatusCode/100 == 4 {
		return newAuthenticationError("Unauthorizatied request", token, nil)
	}

	return newAuthenticationError("Unknown error", token, nil)
}

// addUserInfoToContext adds the information defined in the token to the user context.
//
// We assume that the token in the header has already been checked and is a vaild JWT token.
func addUserInfoToContext(c *context.CustomContext) {
	// See https://github.com/dgrijalva/jwt-go/issues/37 for jwt.Parse with nil
	tokenString := c.Request().Header.Get("Authorization")[7:]
	useTokenToAddUserContext(c, tokenString)
}

func useTokenToAddUserContext(c *context.CustomContext, tokenString string) {
	token, _ := jwt.Parse(tokenString, nil)
	if token == nil {
		panic("Token was not parsable. This should not happen, since we submitted the token to keycloak beforehand.")
	}
	claims := token.Claims.(jwt.MapClaims)

	// Parse roles using a chain of type casts.
	var roles []string
	m1 := claims["realm_access"].(map[string]interface{})
	m2 := m1["roles"].([]interface{})
	for _, v := range m2 {
		roles = append(roles, v.(string))
	}

	// Inject object into context
	auth := context.Authentication{
		Username: claims["preferred_username"].(string),
		Roles:    roles,
	}
	c.Authentication = auth
}

// addEndpoints adds all relevant endpoints for authentication.
func (config *KeycloakConfig) addEndpoints(e *echo.Echo) {
	config.addEndpointLogin(e)
	config.addEndpointRefresh(e)
	config.addEndpointLogout(e)
}

func (config *KeycloakConfig) addEndpointLogin(e *echo.Echo) {
	e.POST(config.LoginURL, func(cc echo.Context) error {
		c := cc.(*context.CustomContext)
		log := c.Log()

		var r loginRequest
		c.Bind(&r)
		if r.Username == "" && r.Password == "" {
			return c.JSON(http.StatusBadRequest, response.NewError("Username and password are empty"))
		}

		token, err := config.handleInitialLogin(c, r)
		if err != nil {
			log.Info().Msg("Login failed: " + err.Error())
			return c.JSON(http.StatusUnauthorized, response.NewError(err.Error()))
		}

		return c.JSON(http.StatusOK, token)
	})
}

func (config *KeycloakConfig) addEndpointRefresh(e *echo.Echo) {
	e.POST(config.RefreshURL, func(cc echo.Context) error {
		c := cc.(*context.CustomContext)
		log := c.Log()

		var r refreshRequest
		c.Bind(&r)
		if r.RefreshToken == "" {
			log.Info().Msg("Empty refresh token")
			return c.JSON(http.StatusBadRequest, response.NewError("Empty refreshToken"))
		}

		token, err := config.handleRefresh(c, r)
		if err != nil {
			log.Info().Msg("Login failed: " + err.Error())
			return c.JSON(http.StatusUnauthorized, response.NewError(err.Error()))
		}

		log = c.Log()
		log.Info().Msg("Successful refresh")
		return c.JSON(http.StatusOK, token)
	})
}

func (config *KeycloakConfig) handleRefresh(c *context.CustomContext, r refreshRequest) (*authenticationResponse, error) {
	// Send request to keycloak
	m := make(map[string][]string)
	m["refresh_token"] = []string{r.RefreshToken}
	m["grant_type"] = []string{"refresh_token"}
	m["client_id"] = []string{"api"}
	resp, err := http.PostForm(config.getKeycloakURLFor("token"), m)
	if err != nil {
		return nil, errors.New("Unknown error")
	}
	log.Info().Int("code", resp.StatusCode).Msg("Status code")
	if resp.StatusCode/100 == 4 {
		return nil, errors.New("Unauthorized")
	}
	dec := json.NewDecoder(resp.Body)
	var v map[string]string
	dec.Decode(&v)

	atoken := v["access_token"]
	token := authenticationResponse{
		AccessToken:  v["access_token"],
		RefreshToken: v["refresh_token"],
	}
	addTokenToCache(atoken)
	useTokenToAddUserContext(c, atoken)
	return &token, nil
}

func (config *KeycloakConfig) handleInitialLogin(c *context.CustomContext, r loginRequest) (*authenticationResponse, error) {
	// Send request to keycloak
	m := make(map[string][]string)
	m["username"] = []string{r.Username}
	m["password"] = []string{r.Password}
	m["grant_type"] = []string{"password"}
	m["client_id"] = []string{"api"}
	resp, err := http.PostForm(config.getKeycloakURLFor("token"), m)
	if err != nil {
		return nil, errors.New("Unknown error")
	}
	log.Info().Int("code", resp.StatusCode).Msg("Status code")
	if resp.StatusCode/100 == 4 {
		return nil, errors.New("Unauthorized")
	}
	dec := json.NewDecoder(resp.Body)
	var v map[string]string
	dec.Decode(&v)

	atoken := v["access_token"]
	token := authenticationResponse{
		AccessToken:  v["access_token"],
		RefreshToken: v["refresh_token"],
	}
	addTokenToCache(atoken)
	useTokenToAddUserContext(c, atoken)
	log := c.Log()
	log.Info().Msg("Successful login")
	return &token, nil
}

func addTokenToCache(token string) {
	t, _ := jwt.Parse(token, nil)
	claims := t.Claims.(jwt.MapClaims)
	expiresAt := int64(claims["exp"].(float64))
	cache[token] = expiresAt
	log.Info().
		Str("token", token).
		Int64("expiresAt", expiresAt).
		Msg("Adding token to cache")
}

func (config *KeycloakConfig) addEndpointLogout(e *echo.Echo) {
	e.POST(config.LogoutURL, func(cc echo.Context) error {
		c := cc.(*context.CustomContext)
		log := c.Log()

		// Try to logout.
		var r logoutRequest
		c.Bind(&r)
		m := make(map[string][]string)
		// TODO Make client_id configurable
		m["client_id"] = []string{"api"}
		m["refresh_token"] = []string{r.RefreshToken}
		m["username"] = []string{r.Username}
		m["password"] = []string{r.Password}
		resp, err := http.PostForm(config.getKeycloakURLFor("logout"), m)
		if err != nil || resp.StatusCode != 204 {
			message, _ := ioutil.ReadAll(resp.Body)
			log.Warn().
				Str("error", string(message)).
				Int("statusCode", resp.StatusCode).
				Msg("Internal server error")
			return c.JSON(http.StatusInternalServerError, response.NewError("Internal error."))
		}

		// Clear cache. This will always work, since we wouldn't be able to call the endpoint
		// without authentication.
		token := c.Request().Header.Get("Authorization")[7:]
		delete(cache, token)

		log.Info().Msg("Logged out successfully")
		return c.NoContent(http.StatusOK)
	})
}
