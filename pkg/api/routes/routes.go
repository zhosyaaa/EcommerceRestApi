package routes

import (
	controllers2 "Ecommerce/pkg/api/controllers"
	middlewares "Ecommerce/pkg/api/middleware"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	controllers2.AdminController
	controllers2.CartController
	controllers2.OrderController
	controllers2.UserController
	controllers2.ProductController
	controllers2.AddressController
}

func NewRoutes(adminController controllers2.AdminController, cartController controllers2.CartController, orderController controllers2.OrderController, userController controllers2.UserController, productController controllers2.ProductController, addressController controllers2.AddressController) *Routes {
	return &Routes{AdminController: adminController, CartController: cartController, OrderController: orderController, UserController: userController, ProductController: productController, AddressController: addressController}
}

func (r *Routes) SetupRoutes(app *gin.Engine) *gin.Engine {
	api := app.Group("api/v1")

	//auth routes
	userApi := api.Group("/users/auth")
	{
		userApi.POST("/singup", r.Signup)
		userApi.POST("/singin", r.Signin)
		userApi.POST("/singout", middlewares.RequireAuthMiddleware, r.Signout)
		userApi.GET("/profile", middlewares.RequireAuthMiddleware, r.Profile)
	}

	// products routes
	productsApi := api.Group("/products")
	{
		productsApi.GET("/", r.GetAllProducts)
		productsApi.GET("/:id", r.GetProduct)
		productsApi.POST("/create", middlewares.RequireAuthMiddleware, r.CreateProduct)
		productsApi.PUT("/:id", middlewares.RequireAuthMiddleware, r.UpdateProduct)
		productsApi.DELETE("/:id", middlewares.RequireAuthMiddleware, r.DeleteProduct)
	}

	// address routes
	addressApi := api.Group("/address", middlewares.RequireAuthMiddleware)
	{
		addressApi.PUT("/update/:id", r.UpdateAddress)
	}

	// cart routes
	cartApi := api.Group("/cart")
	{
		cartApi.POST("/remove/:id", middlewares.RequireAuthMiddleware, r.RemoveProductFromCart)
		cartApi.POST("/add/:id", middlewares.RequireAuthMiddleware, r.AddProductToCart)
	}
	// order routes
	orderApi := api.Group("order", middlewares.RequireAuthMiddleware)
	{
		orderApi.POST("/", r.OrderAll)
		orderApi.POST("/:id", r.OrderOne)
	}

	// admin routes
	adminApi := api.Group("admin", middlewares.RequireAuthMiddleware)
	{
		adminApi.GET("/getUser/:id", r.GetUser)
		adminApi.GET("getUsers", r.GetUsers)
		adminApi.DELETE("/deleteUser/:id", r.DeleteUser)
		adminApi.DELETE("/deleteUsers", r.DeleteAllUsers)
	}
	return app
}
