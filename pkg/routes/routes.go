package routes

import (
	controllers "Ecommerce/pkg/controllers"
	"github.com/gin-gonic/gin"
)

func Routes(app *gin.Engine) *gin.Engine {
	api := app.Group("api/v1")

	//auth routes
	userApi := api.Group("/users/auth")
	{
		userApi.POST("singup", controllers.Signup)   //, middlewares.ValidateCredentialsMiddleware +
		userApi.POST("singin", controllers.Signin)   //middlewares.ValidateCredentialsMiddleware, +
		userApi.POST("singout", controllers.Signout) //  middlewares.RequireAuthMiddleware,
		userApi.GET("/profile", controllers.Profile) //  middlewares.RequireAuthMiddleware,
	}

	// products routes
	productsApi := api.Group("/products")
	{
		productsApi.GET("/", controllers.GetAllProducts)
		productsApi.GET("/:id", controllers.GetProduct)
		productsApi.POST("/create", controllers.CreateProduct) //  middlewares.RequireAuthMiddleware,
		productsApi.PUT("/:id", controllers.UpdateProduct)     //  middlewares.RequireAuthMiddleware,
		productsApi.DELETE("/:id", controllers.DeleteProduct)  ///  middlewares.RequireAuthMiddleware,
	}

	// cart routes
	cartApi := api.Group("/cart") //  middlewares.RequireAuthMiddleware
	{
		cartApi.POST("/remove/:id", controllers.RemoveProductFromCart)
		cartApi.POST("/add/:id", controllers.AddProductToCart)
	}

	// address routes
	addressApi := api.Group("/address") // middlewares.RequireAuthMiddleware
	{
		addressApi.PUT("/update/:id", controllers.UpdateAddress)
	}

	// order routes
	orderApi := api.Group("order") // middlewares.RequireAuthMiddleware
	{
		orderApi.POST("/", controllers.OrderAll)
		orderApi.POST("/:id", controllers.OrderOne)
	}

	// admin routes
	adminApi := api.Group("admin") //  middlewares.RequireAuthMiddleware
	{
		adminApi.GET("/getUser/:id", controllers.GetUser)
		adminApi.GET("getUsers", controllers.GetUsers)
		adminApi.DELETE("/deleteUser/:id", controllers.DeleteUser)
		adminApi.DELETE("/deleteUsers", controllers.DeleteAllUsers)
	}
	return app
}
