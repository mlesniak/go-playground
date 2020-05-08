package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	// RefreshURL string
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type logoutRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}

type authenticationResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// KeycloakWithConfig registers authentication endpoints and handles token validation on each request.
func KeycloakWithConfig(e *echo.Echo, config KeycloakConfig) func(next echo.HandlerFunc) echo.HandlerFunc {
	if config.Protocol == "" || config.Hostname == "" || config.Port == "" || config.Realm == "" || config.LoginURL == "" || config.LogoutURL == "" {
		panic("The keycloak configuration is invalid, at least one required property is empty.")
	}
	config.addEndpoints(e)

	// Check token on each request.
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
				if err := config.isAuthenticated(c); err != nil {
					e := err.(Error)
					log.Info().Str("token", e.Token).Msg(e.Text)
					return c.JSON(http.StatusUnauthorized, response.NewError("Token is invalid"))
				}
				// Everything is ok, add authentication info to request context.
				addUserInfoToContext(c)
				cc := c.(*context.CustomContext)
				log.Info().Str("username", cc.Username()).Msg("Successful authentication by token")
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
func (config *KeycloakConfig) isAuthenticated(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return NewError(nil, "Empty token", token)
	}

	// Access userinfo in keycloak using the provided token to check if the token is valid.
	req, err := http.NewRequest("GET", config.getKeycloakURLFor("userinfo"), nil)
	req.Header.Add("Authorization", token)
	cl := &http.Client{}
	resp, err := cl.Do(req)

	if err != nil {
		return NewError(nil, "Unable to get userinfo", token)
	}
	if resp.StatusCode == 200 {
		return nil
	}
	if resp.StatusCode%100 == 4 {
		return NewError(nil, "Unauthorizatied request", token)
	}

	return NewError(nil, "Unknown error", token)
}

// addUserInfoToContext adds the information defined in the token to the user context.
//
// We assume that the token in the header has already been checked and is a vaild JWT token.
func addUserInfoToContext(cc echo.Context) {
	c := cc.(*context.CustomContext)

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

		var r loginRequest
		c.Bind(&r)
		if r.Username == "" && r.Password == "" {
			return c.JSON(http.StatusBadRequest, response.NewError("Username and password are empty"))
		}

		return config.handleInitialLogin(c, r)
	})
}

func (config *KeycloakConfig) handleInitialLogin(c *context.CustomContext, r loginRequest) error {
	// Send request to keycloak
	m := make(map[string][]string)
	m["username"] = []string{r.Username}
	m["password"] = []string{r.Password}
	m["grant_type"] = []string{"password"}
	m["client_id"] = []string{"api"}
	resp, err := http.PostForm(config.getKeycloakURLFor("token"), m)
	if err != nil {
		return NewError(err, "Unknown error", "")
	}
	if resp.StatusCode%100 == 4 {
		return NewError(err, "Unauthorized", "")
	}
	dec := json.NewDecoder(resp.Body)
	var v map[string]string
	dec.Decode(&v)

	token := authenticationResponse{
		AccessToken:  v["access_token"],
		RefreshToken: v["refresh_token"],
	}
	log.Info().Str("username", r.Username).Msg("Successful login")
	return c.JSON(http.StatusOK, token)
}

func (config *KeycloakConfig) addEndpointLogout(e *echo.Echo) {
	e.POST(config.LogoutURL, func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.NoContent(http.StatusOK)
		}

		var r logoutRequest
		c.Bind(&r)
		m := make(map[string][]string)
		m["client_id"] = []string{"api"}
		m["refresh_token"] = []string{r.RefreshToken}
		m["username"] = []string{r.Username}
		m["password"] = []string{r.Password}
		resp, err := http.PostForm(config.getKeycloakURLFor("logout"), m)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		if resp.StatusCode == 204 {
			return c.NoContent(http.StatusOK)
		}

		return c.NoContent(http.StatusInternalServerError)
	})
}
