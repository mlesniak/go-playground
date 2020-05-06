package authentication

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/mlesniak/go-demo/pkg/context"
	"github.com/mlesniak/go-demo/pkg/errors"
	logger "github.com/mlesniak/go-demo/pkg/log"
)

var log = logger.New()

// KeycloakConfig defines configuration options for the middleware.
type KeycloakConfig struct {
	IgnoredURL []string
	// Add configuration for Keycloak here...
}

// ContextName is the name of the stored authentication object in the context object.
// TODO Use custom context here.
// TODO Add methods to authentication such as hasRole(string) bool
const ContextName = "Authentication"

// Authentication contains all user information of an authenticated user.
type Authentication struct {
	Username string
	Roles    []string
}

// KeycloakWithConfig ... with config
func KeycloakWithConfig(config KeycloakConfig) func(next echo.HandlerFunc) echo.HandlerFunc {
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
				if !IsAuthenticated(c) {
					return c.JSON(http.StatusUnauthorized, errors.Response{"Token is invalid"})
				}
				addUserInfoToContext(c)
			}

			// Continue our chain.
			if err := next(c); err != nil {
				c.Error(err)
			}
			return nil
		}
	}
}

// IsAuthenticated returns true if the user submits a valid JWT token.
func IsAuthenticated(c echo.Context) bool {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		log.Info("No token provided in Authorization header")
		return false
	}

	req, err := http.NewRequest("GET", "http://localhost:8081/auth/realms/mlesniak/protocol/openid-connect/userinfo", nil)
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
func addUserInfoToContext(c echo.Context) {
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
	auth := Authentication{
		Username: claims["preferred_username"].(string),
		Roles:    roles,
	}
	c.Set(ContextName, auth)
}

// AddAuthenticationEndpoints adds the login endpoint for authentication.
// WILL BE HEAVILY REFACTORED INTO ITS OWN PACKAGE...
func AddAuthenticationEndpoints(e *echo.Echo) {
	type request struct {
		Username     string `json:"username"`
		Password     string `json:"password"`
		RefreshToken string `json:"refreshToken"`
	}

	e.POST("/api/logout", func(c echo.Context) error {
		log.Info("/api/logout called")
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			log.Info("/logout called for non-authorized user")
			return c.NoContent(http.StatusOK)
		}

		// c.Request().ParseForm()
		// r := c.Request().Form.Get("request_token")
		var r request
		c.Bind(&r)
		m := make(map[string][]string)
		m["client_id"] = []string{"api"}
		m["refresh_token"] = []string{r.RefreshToken}
		m["username"] = []string{r.Username}
		m["password"] = []string{r.Password}
		resp, err := http.PostForm("http://localhost:8081/auth/realms/mlesniak/protocol/openid-connect/logout", m)
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

	// This endpoint can be used for both first and later authentication with refresh tokens.
	e.POST("/api/login", func(cc echo.Context) error {
		c := cc.(*context.CustomContext)
		c.Foo()

		type response struct {
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
		}

		var r request
		c.Bind(&r)

		if r.Username != "" && r.Password != "" {
			log.WithField("username", r.Username).Info("Login attempt")
			// Send request to keycloak
			m := make(map[string][]string)
			m["username"] = []string{r.Username}
			m["password"] = []string{r.Password}
			m["grant_type"] = []string{"password"}
			m["client_id"] = []string{"api"}
			resp, err := http.PostForm("http://localhost:8081/auth/realms/mlesniak/protocol/openid-connect/token", m)
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

		return c.String(http.StatusOK, "/api/login case not implemented")
	})
}
