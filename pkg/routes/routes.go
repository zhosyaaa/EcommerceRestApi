package routes

import (
	middlewares "Ecommerce/pkg/middleware"
	controllers "Ecommerce/pkg/controllers"
	"github.com/gin-gonic/gin"
)

func Routes(app *gin.Engine) {
	api := app.Group("api/v1")

	//auth routes
	userApi := api.Group("/users/auth")
	{
		userApi.POST("singup", middlewares.ValidateCredentialsMiddleware, controllers.Signup)
		userApi.POST("singin", middlewares.ValidateCredentialsMiddleware, controllers.Signin)
		userApi.POST("singout", middlewares.RequireAuthMiddleware, controllers.Signout)
		userApi.GET("/profile", middlewares.RequireAuthMiddleware, controllers.Profile)
	}

	// products routes
	productsApi := api.Group("/products")
	{
		productsApi.GET("/", controllers.GetAllProducts)
		productsApi.GET("/:id", controllers.GetProduct)
		productsApi.POST("/create", middlewares.RequireAuthMiddleware, controllers.CreateProduct)
		productsApi.PUT("/:id", middlewares.RequireAuthMiddleware, controllers.UpdateProduct
		productsApi.DELETE("/:id", middlewares.RequireAuthMiddleware, controllers.DeleteProduct)
	}

	// cart routes
	cartApi := api.Group("/cart",  middlewares.RequireAuthMiddleware)
	{
		cartApi.POST("/remove/:id",  controllers.RemoveProductFromCart)
		cartApi.POST("/add/:id",  controllers.AddProductToCart)
	}

	// address routes
	addressApi := api.Group("/address", middlewares.RequireAuthMiddleware)
	{
		addressApi.PUT("/update",controllers.UpdateAddress)
	}

	// order routes
	orderApi := api.Group("order", middlewares.RequireAuthMiddleware)
	{
		orderApi.POST("/", controllers.OrderAll)
		orderApi.POST("/:id",  controllers.OrderOne)
	}

	// admin routes
	adminApi := api.Group("admin", middlewares.RequireAuthMiddleware)
	{
		adminApi.GET("/getUser/:id",  controllers.GetUser)
		adminApi.GET("getUsers", controllers.GetUsers)
		adminApi.DELETE("/deleteUser/:id", controllers.DeleteUser)
		adminApi.DELETE("/deleteUsers", controllers.DeleteAllUsers)
	}
}
