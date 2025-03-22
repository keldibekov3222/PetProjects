package routes

import (
	"github.com/gin-gonic/gin"
	"order-service/handlers"
)

func RegisterRoutes(r *gin.Engine, userHandler *handlers.UserHandler, orderHandler *handlers.OrderHandler, productHandler *handlers.ProductHandler) {
	// Регистрация маршрутов для пользователей
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.GET("/users", userHandler.GetAllUsers)
	r.GET("/users/:id", userHandler.GetUserByID)
	r.PUT("/users/:id", userHandler.UpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)

	// Регистрация маршрутов для заказов
	r.POST("/orders", orderHandler.CreateOrder)
	r.GET("/orders/:id", orderHandler.GetOrderById)
	r.GET("/orders/", orderHandler.GetAllOrders)
	r.DELETE("/orders/:id", orderHandler.DeleteOrder)
	r.PUT("/orders/:id", orderHandler.UpdateOrder)

	// Регистрация маршрутов для продуктов
	r.POST("/products", productHandler.CreateProduct)
	r.GET("/products", productHandler.GetAllProducts)
	r.GET("/products/:id", productHandler.GetProductById)
	r.PUT("/products/:id", productHandler.UpdateProduct)
	r.DELETE("/products/:id", productHandler.DeleteProduct)
}
