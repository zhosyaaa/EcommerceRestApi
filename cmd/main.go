package main

import (
	controllers2 "Ecommerce/pkg/api/controllers"
	middleware2 "Ecommerce/pkg/api/middleware"
	"Ecommerce/pkg/api/routes"
	"Ecommerce/pkg/config"
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/repository"
	interfaces "Ecommerce/pkg/repository/interface"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var userRepository interfaces.UserRepository
var productRepository interfaces.ProductRepository
var orderRepository interfaces.OrderRepository
var addressRepository interfaces.AddressRepository

var adminController controllers2.AdminController
var cartController controllers2.CartController
var orderController controllers2.OrderController
var userController controllers2.UserController
var productController controllers2.ProductController
var addressController controllers2.AddressController

func main() {
	configs, err := config.LoadConfig()
	if err != nil {
		fmt.Println("loadConfig error")
	}
	db.ConnectDB(configs)
	app := gin.New()
	app.Use(gin.Recovery())
	_ = app.SetTrustedProxies(nil)
	app.Use(middleware2.CORSMiddleware())

	userRepository = repository.NewUserRepository(db.GetDB())
	productRepository = repository.NewProductRepository(db.GetDB())
	orderRepository = repository.NewOrderRepository(db.GetDB())
	addressRepository = repository.NewAddressDatabase(db.GetDB())

	adminController := controllers2.NewAdminController(userRepository)
	cartController := controllers2.NewCartController(orderRepository, productRepository, userRepository)
	orderController := controllers2.NewOrderController(orderRepository, productRepository, userRepository)
	userController := controllers2.NewUserController(userRepository)
	productController := controllers2.NewProductController(productRepository)
	addressController := controllers2.NewAddressController(addressRepository, userRepository)

	routeHandlers := routes.NewRoutes(*adminController, *cartController, *orderController, *userController, *productController, *addressController)

	app = routeHandlers.SetupRoutes(app)
	addr := config.GetEnvVar("GIN_ADDR")
	port := config.GetEnvVar("GIN_PORT")

	https := config.GetEnvVar("GIN_HTTPS")
	if https == "true" {
		certFile := config.GetEnvVar("GIN_CERT")
		certKey := config.GetEnvVar("GIN_CERT_KEY")

		if err := app.RunTLS(fmt.Sprintf("%s:%s", addr, port), certFile, certKey); err != nil {
			log.Fatal().Err(err).Msg("Error occurred while setting up the server in HTTPS mode")
		}
	}

	log.Info().Msgf("Starting service on http//:%s:%s", addr, port)
	if err := app.Run(fmt.Sprintf("%s:%s", addr, port)); err != nil {
		log.Fatal().Err(err).Msg("Error occurred while setting up the server")
	}
}
