package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	--"github.com/gofiber/fiber/v2/routes"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
//	routes.Router(app)
}
