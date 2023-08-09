package main

import (
	"Ecommerce/pkg/config"
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/routes"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println("loadConfig error")
	}
	db.ConnectDB(config)
	app := gin.New()
	app.Use(gin.Logger())

	app = routes.Routes(app)

	app.Run(":8080")
}
