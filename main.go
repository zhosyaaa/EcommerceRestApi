package main

import (
	"Ecommerce/pkg/config"
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/middleware"
	"Ecommerce/pkg/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Setting configs")
	configs, err := config.LoadConfig()
	if err != nil {
		fmt.Println("loadConfig error")
	}

	log.Info().Msg("Connect with database")
	db.ConnectDB(configs)

	log.Info().Msg("Initializing service")
	app := gin.New()
	app.Use(gin.Recovery())
	_ = app.SetTrustedProxies(nil)

	log.Info().Msg("Adding cors and request logging middleware")
	app.Use(middleware.CORSMiddleware(), middleware.RequestID(), middleware.RequestLogger())

	log.Info().Msg("Setting up routers")
	app = routes.Routes(app)

	addr := config.GetEnvVar("GIN_ADDR")
	port := config.GetEnvVar("GIN_PORT")

	https := config.GetEnvVar("GIN_HTTPS")
	if https == "true" {
		certFile := config.GetEnvVar("GIN_CERT")
		certKey := config.GetEnvVar("GIN_CERT_KEY")
		log.Info().Msgf("Starting service on https//:%s:%s", addr, port)

		if err := app.RunTLS(fmt.Sprintf("%s:%s", addr, port), certFile, certKey); err != nil {
			log.Fatal().Err(err).Msg("Error occurred while setting up the server in HTTPS mode")
		}
	}

	log.Info().Msgf("Starting service on http//:%s:%s", addr, port)
	if err := app.Run(fmt.Sprintf("%s:%s", addr, port)); err != nil {
		log.Fatal().Err(err).Msg("Error occurred while setting up the server")
	}
}
