package main

import (
	"context"
	"log"
	"order-service/config"
	"order-service/handlers"
	"order-service/repositories"
	"order-service/routes"
	"order-service/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Подключение к PostgreSQL для заказов
	db, err := repositories.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Println("Closing database connection")
		db.Close(context.Background())
	}()

	// Подключение к MongoDB для продуктов и пользователей
	mongoRepo, err := repositories.NewMongoDBRepository(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", err)
	}

	// Подключение к Redis для корзины
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	// Репозитории
	orderRepo := repositories.NewOrderRepository(db)
	productRepo := repositories.NewProductRepository(mongoRepo.DB)
	userRepo := repositories.NewUserRepository(mongoRepo.DB)

	// Сервисы
	orderService := services.NewOrderService(orderRepo, redisClient)
	productService := services.NewProductService(productRepo, redisClient)
	userService := services.NewUserService(userRepo, redisClient)
	cartService := services.NewCartService(redisClient, productRepo, orderRepo)

	// Хендлеры
	orderHandler := handlers.NewOrderHandler(orderService)
	productHandler := handlers.NewProductHandler(productService)
	userHandler := handlers.NewUserHandler(userService)
	cartHandler := handlers.NewCartHandler(cartService)

	// Создание и настройка Gin
	r := gin.Default()

	// Регистрация маршрутов
	routes.RegisterRoutes(r, userHandler, orderHandler, productHandler, cartHandler)

	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
