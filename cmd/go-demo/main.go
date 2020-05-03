package main

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	logger "github.com/mlesniak/go-demo/pkg/log"
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
	addVersionEndpoint(e)
	addAPIEndpoint(e)
	addAuthenticationEndpoints(e)

	log.Info("Application started")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Info(e.Start(":" + port))
}

func addVersionEndpoint(e *echo.Echo) {
	e.GET("/api/version", func(c echo.Context) error {
		log := logger.AddUser(log, c)

		log.Info("Version info requested")
		commit := os.Getenv("COMMIT")
		if commit == "" {
			commit = "<No COMMIT environment variable set>"
		}
		return c.JSON(http.StatusOK, struct {
			Version string `json:"version"`
		}{commit})
	})
}

func addAPIEndpoint(e *echo.Echo) {
	type request struct {
		Number int `json:"number"`
	}

	type response struct {
		Number int `json:"number"`
	}

	e.POST("/api", func(c echo.Context) error {
		log := logger.AddUser(log, c)

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

// Will be refactored into its own package...
func addAuthenticationEndpoints(e *echo.Echo) {
	// This endpoint can be used for both first and later authentication with refresh tokens.
	e.POST("/api/login", func(c echo.Context) error {
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
			bs, _ := ioutil.ReadAll(resp.Body)
			log.Info("Response ", string(bs))
			// dec := json.NewDecoder(resp.Body)
			// var v interface{}
			// dec.Decode(&v)
			// log.Info("Response {}", v)
		}

		return c.String(http.StatusOK, "/api/login working")
	})
}

// else if r.RefreshToken != "" {
// } else {
// 	return c.JSON(http.StatusBadRequest, ErrorResponse{"Missing parameters"})
// }
