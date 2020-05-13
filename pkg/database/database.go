package database

import (
	ctx "context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// Get returns the default collection for this application.
func Get() *mongo.Collection {
	api := client.Database("dev").Collection("api")
	return api
}

// DefaultContext returns a context with a default timeout.
// Note that we ignore the cancel function on purpose for now.
func DefaultContext() ctx.Context {
	timeout := time.Second * 3
	c, _ := ctx.WithTimeout(ctx.Background(), timeout)
	return c
}

func init() {
	c, cancel := ctx.WithTimeout(ctx.Background(), 10*time.Second)
	defer cancel()

	// TODO Use environment variables
	host := "mongodb://localhost:27017"
	options := options.Client().
		SetAuth(options.Credential{
			Username: "admin",
			Password: "admin",
		}).
		ApplyURI(host)
	cl, err := mongo.Connect(c, options)
	if err != nil {
		log.Error().Msg("Unable to initialize database")
		panic(err)
	}
	client = cl
	log.Info().
		Str("host", host).
		Msg("Database initialized")
}
