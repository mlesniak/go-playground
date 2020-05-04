package authentication

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	logger "github.com/mlesniak/go-demo/pkg/log"
)

var log = logger.New()

// IsAuthenticated returns true if the user submits a valid JWT token.
func IsAuthenticated(c echo.Context) bool {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return false
	}
	log.Info("token:", token)

	// Access keycloak until we use caching.
	// m := make(map[string][]string)
	// m["refresh_token"] = ...
	// m["grant_type"] = []string{"password"}
	// m["client_id"] = []string{"api"}

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

// else if r.RefreshToken != "" {
// } else {
// 	return c.JSON(http.StatusBadRequest, ErrorResponse{"Missing parameters"})
// }

// AddAuthenticationEndpoints adds the login endpoint for authentication.
// WILL BE HEAVILY REFACTORED INTO ITS OWN PACKAGE...
func AddAuthenticationEndpoints(e *echo.Echo) {
	// This endpoint can be used for both first and later authentication with refresh tokens.
	e.POST("/api/login", func(c echo.Context) error {
		type response struct {
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
		}

		type request struct {
			Username     string `json:"username"`
			Password     string `json:"password"`
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
			log.Info("Response {}", v)

			token := response{
				AccessToken:  v["access_token"],
				RefreshToken: v["refresh_token"],
			}
			return c.JSON(http.StatusOK, token)
		}

		return c.String(http.StatusOK, "/api/login case not implemented")
	})
}