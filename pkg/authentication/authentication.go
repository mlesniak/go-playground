package authentication

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	logger "github.com/mlesniak/go-demo/pkg/log"
)

var log = logger.New()

// KeycloakMiddleware defines a middleware to check authentication.
func KeycloakMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		check := true
		ignoreURL := []string{
			"/api/login",
		}
		log.Info("Request from ", c.Request().RequestURI)
		for _, i := range ignoreURL {
			if i == c.Request().URL.RequestURI() {
				check = false
				break
			}
		}

		if check && !IsAuthenticated(c) {
			return c.NoContent(http.StatusUnauthorized)
		}

		// TODO If authenticated, add user and roles to request context.

		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}


// IsAuthenticated returns true if the user submits a valid JWT token.
func IsAuthenticated(c echo.Context) bool {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		log.Info("No token")
		return false
	}

	log.Info("Check token", token)
	req, err := http.NewRequest("GET", "http://localhost:8081/auth/realms/mlesniak/protocol/openid-connect/userinfo", nil)
	req.Header.Add("Authorization", token)
	cl := &http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		return false
	}
	if resp.StatusCode == 200 {
		return true
	}

	return false
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
	e.POST("/api/login", func(c echo.Context) error {
		type response struct {
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
		}

		var r request
		c.Bind(&r)
		log.Info("Body:", r)

		if r.Username != "" && r.Password != "" {
			log.WithField("username", r.Username).Info("Login")
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
			log.WithField("code", resp.StatusCode).Info("Successful login")
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
