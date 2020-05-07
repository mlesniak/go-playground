package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/mlesniak/go-demo/pkg/context"
	"github.com/mlesniak/go-demo/pkg/errors"
	logger "github.com/mlesniak/go-demo/pkg/log"
)

// TODO change log library to zerolog.
var log = logger.New()

// KeycloakConfig defines configuration options for the middleware.
type KeycloakConfig struct {
	Protocol   string
	Hostname   string
	Port       string
	Realm      string

	IgnoredURL []string
	LoginURL   string
	LogoutURL  string
	// RefreshURL string
}

// request is the default request for login and refresh attempts.
// TODO HAVING A SINGLE ENDPOINT FOR LOGIN AND REFRESH IS A STUPID IDEA. CHANGE THIS.
type request struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	RefreshToken string `json:"refreshToken"`
}

type response struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// KeycloakWithConfig ... with config
func KeycloakWithConfig(e *echo.Echo, config KeycloakConfig) func(next echo.HandlerFunc) echo.HandlerFunc {
	if config.Protocol == "" || config.Hostname == "" || config.Port == "" || config.Realm == "" || config.LoginURL == "" || config.LogoutURL == "" {
		panic("The keycloak configuration is invalid, at least one required property is empty.")
	}

	config.addEndpoints(e)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
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
				if !config.isAuthenticated(c) {
					return c.JSON(http.StatusUnauthorized, errors.Response{"Token is invalid"})
				}
				// Everything is ok, add authentication info to request context.
				addUserInfoToContext(c)
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
// TODO Add caching while respecting the expiration date.
func (config *KeycloakConfig) isAuthenticated(c echo.Context) bool {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		log.Info("No token provided in Authorization header")
		return false
	}

	// Access userinfo in keycloak using the provided token to check if the token is valid.
	req, err := http.NewRequest("GET", config.getKeycloakURLFor("userinfo"), nil)
	req.Header.Add("Authorization", token)
	cl := &http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		log.Info("Unable to authorize with token")
		return false
	}
	if resp.StatusCode == 200 {
		return true
	}
	return false
}

// addUserInfoToContext adds the information defined in the token to the user context.
func addUserInfoToContext(cc echo.Context) {
	c := cc.(*context.CustomContext)

	// If authenticated, add username and roles to request context for later processing.
	// See https://github.com/dgrijalva/jwt-go/issues/37 for jwt.Parse with nil
	tokenString := c.Request().Header.Get("Authorization")[7:]
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
	config.addEndpointLogout(e)
}

func (config *KeycloakConfig) addEndpointLogin(e *echo.Echo) {
	e.POST(config.LoginURL, func(cc echo.Context) error {
		c := cc.(*context.CustomContext)
		c.Username()

		var r request
		c.Bind(&r)

		if r.Username == "" && r.Password == "" && r.RefreshToken == "" {
			return c.NoContent(http.StatusBadRequest)
		}

		// Case: initial login, i.e. no refresh token submitted.
		if r.Username != "" && r.Password != "" && r.RefreshToken == "" {
			return config.handleInitialLogin(c, r)
		}

		// TODO implement refresh token

		return c.String(http.StatusOK, "/api/login case not implemented")
	})
}

func (config *KeycloakConfig) handleInitialLogin(c *context.CustomContext, r request) error {
	log.WithField("username", r.Username).Info("Login attempt")
	// Send request to keycloak
	m := make(map[string][]string)
	m["username"] = []string{r.Username}
	m["password"] = []string{r.Password}
	m["grant_type"] = []string{"password"}
	m["client_id"] = []string{"api"}
	resp, err := http.PostForm(config.getKeycloakURLFor("token"), m)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic("Not working:" + string(resp.StatusCode))
	}
	log.WithField("username", r.Username).Info("Login successful")
	dec := json.NewDecoder(resp.Body)
	var v map[string]string
	dec.Decode(&v)

	token := response{
		AccessToken:  v["access_token"],
		RefreshToken: v["refresh_token"],
	}
	return c.JSON(http.StatusOK, token)
}

func (config *KeycloakConfig) addEndpointLogout(e *echo.Echo) {
	e.POST(config.LogoutURL, func(c echo.Context) error {
		log.Info("/api/logout called")
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			log.Info("/logout called for non-authorized user")
			return c.NoContent(http.StatusOK)
		}

		var r request
		c.Bind(&r)
		m := make(map[string][]string)
		m["client_id"] = []string{"api"}
		m["refresh_token"] = []string{r.RefreshToken}
		m["username"] = []string{r.Username}
		m["password"] = []string{r.Password}
		resp, err := http.PostForm(config.getKeycloakURLFor("logout"), m)
		if err != nil {
			log.Warn(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		if resp.StatusCode == 204 {
			log.Info("Logout successful")
			return c.NoContent(http.StatusOK)
		}

		log.Info("Oops: ", resp.StatusCode)
		return c.NoContent(http.StatusInternalServerError)
	})
}
