package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	logger "github.com/mlesniak/go-demo/pkg/log"
	"github.com/mlesniak/go-demo/pkg/version"
)

var log = logger.New()

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Middlewares.
	// Disabled until we add keycloak handling.
	// addJWTMiddleware(e)

	// Endpoints.
	version.AddVersionEndpoint(e)
	addAPIEndpoint(e)
	addAuthenticationEndpoints(e)

	log.Info("Application started")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Info(e.Start(":" + port))
}

func addAPIEndpoint(e *echo.Echo) {
	type request struct {
		Number int `json:"number"`
	}

	type response struct {
		Number int `json:"number"`
	}

	e.POST("/api", func(c echo.Context) error {
		// log := logger.AddUser(log, c)
		// This will be done as a middleware, later on.
		authenticated := isAuthenticated(c)
		if !authenticated {
			return c.NoContent(http.StatusUnauthorized)
		}

		var json request
		err := c.Bind(&json)
		if err != nil {
			log.WithField("error", err).Warn("Unable to parse json")
			return c.String(http.StatusBadRequest, "Unable to parse request")
		}
		log.WithField("number", json.Number).Info("Request received")
		resp := response{json.Number + 1}
		return c.JSON(http.StatusOK, resp)
	})
}

func addJWTMiddleware(e *echo.Echo) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("NO JWT_SECRET set. Aborting.")
	}
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(secret),
		Skipper: func(c echo.Context) bool {
			// List of urls to ignore for authentication.
			ignoredURL := []string{
				"/api/version",
			}

			path := c.Request().URL.Path
			for _, v := range ignoredURL {
				if v == path {
					return true
				}
			}
			return false
		},
	}))
}

// WILL BE HEAVILY REFACTORED INTO ITS OWN PACKAGE...
func addAuthenticationEndpoints(e *echo.Echo) {
	// This endpoint can be used for both first and later authentication with refresh tokens.
	e.POST("/api/login", func(c echo.Context) error {
		type response struct {
			AccessToken     string `json:"accessToken"`
			RefreshToken     string `json:"refreshToken"`
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
				AccessToken: v["access_token"],
				RefreshToken: v["refresh_token"],
			}
			return c.JSON(http.StatusOK, token)
		}

		return c.String(http.StatusOK, "/api/login case not implemented")
	})
}

// Returns false if unauthorized.
func isAuthenticated(c echo.Context) bool {
	bearer := c.Request().Header.Get("Authorization")
	if bearer == "" {
		return false
	}

	return true
}

// else if r.RefreshToken != "" {
// } else {
// 	return c.JSON(http.StatusBadRequest, ErrorResponse{"Missing parameters"})
// }
