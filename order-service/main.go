package main

import (
	"context"
	"flag"
	"log"
	"order-service/config"
	"order-service/db"
	"order-service/handlers"
	"order-service/repositories"
	"order-service/routes"
	"order-service/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Определяем флаг для запуска только миграций
	migrateOnly := flag.Bool("migrate", false, "Run database migrations only")
	flag.Parse()

	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Подключение к PostgreSQL для заказов и пользователей
	dbConn, err := repositories.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Println("Closing database connection")
		dbConn.Close(context.Background())
	}()

	// Применяем миграции
	if err := db.MigrateConfig(cfg); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Если указан флаг -migrate, завершаем работу после миграций
	if *migrateOnly {
		log.Println("Migrations completed successfully. Exiting.")
		return
	}

	// Подключение к MongoDB для продуктов
	mongoRepo, err := repositories.NewMongoDBRepository(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", err)
	}

	// Подключение к Redis для корзины
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	// Репозитории
	orderRepo := repositories.NewOrderRepository(dbConn)
	userRepo := repositories.NewUserRepository(dbConn)
	productRepo := repositories.NewProductRepository(mongoRepo.DB)

	// Сервисы
	orderService := services.NewOrderService(orderRepo, redisClient)
	userService := services.NewUserService(userRepo, redisClient)
	productService := services.NewProductService(productRepo, redisClient)
	cartService := services.NewCartService(redisClient, productRepo, orderRepo, userRepo)

	// Хендлеры
	orderHandler := handlers.NewOrderHandler(orderService)
	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)
	cartHandler := handlers.NewCartHandler(cartService)

	// Создание и настройка Gin
	r := gin.Default()

	// Регистрация маршрутов
	routes.RegisterRoutes(r, userHandler, orderHandler, productHandler, cartHandler)

	// Запуск сервера
	serverAddr := cfg.ServerAddr
	if serverAddr == "" {
		serverAddr = ":8080" // Default address if not configured
	}
	r.Run(serverAddr)
}
