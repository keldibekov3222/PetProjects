package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"order-service/config"
	"order-service/handlers"
	"order-service/repositories"
	"order-service/services"
)

func main() {
	cfg := config.LoadConfig()
	db, err := repositories.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Println("Closing database connection")
		db.Close(context.Background())
	}()

	mongoRepo, err := repositories.NewMongoDBRepository(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", err)
	}

	orderRepo := repositories.NewOrderRepository(db)
	if orderRepo == nil {
		log.Fatal("orderRepo is nil")
	}
	productRepo := repositories.NewProductRepository(mongoRepo.DB)

	orderService := services.NewOrderService(orderRepo)
	productService := services.NewProductService(productRepo)

	orderHandler := handlers.NewOrderHandler(orderService)
	pruductHandler := handlers.NewProductHandler(productService)

	r := gin.Default()

	r.POST("/orders", orderHandler.CreateOrder)
	r.GET("/orders/:id", orderHandler.GetOrderById)
	r.GET("/orders/", orderHandler.GetAllOrders)
	r.DELETE("/orders/:id", orderHandler.DeleteOrder)
	r.PUT("/orders/:id", orderHandler.UpdateOrder)

	r.POST("/products", pruductHandler.CreateProduct)
	r.GET("/products", pruductHandler.GetAllProducts)
	r.GET("/products/:id", pruductHandler.GetProductById)
	r.PUT("/products/:id", pruductHandler.UpdateProduct)
	r.DELETE("/products/:id", pruductHandler.DeleteProduct)

	r.Run(":8080")
}
