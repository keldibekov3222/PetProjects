package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"order-service/config"
	"order-service/handlers"
	"order-service/repositories"
	"order-service/routes"
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
	userRepo := repositories.NewUserRepository(mongoRepo.DB)

	orderService := services.NewOrderService(orderRepo)
	productService := services.NewProductService(productRepo)
	userService := services.NewUserService(userRepo)

	orderHandler := handlers.NewOrderHandler(orderService)
	productHandler := handlers.NewProductHandler(productService)
	userHandler := handlers.NewUserHandler(userService)

	r := gin.Default()

	// Вызов функции для регистрации маршрутов
	routes.RegisterRoutes(r, userHandler, orderHandler, productHandler)

	r.Run(":8080")
}
