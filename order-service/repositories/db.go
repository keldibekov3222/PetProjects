package repositories

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
	"log"
	"order-service/config"
)

func ConnectDB(cfg *config.Config) (*pgx.Conn, error) {
	// Формируем строку подключения (DSN)
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	log.Println("Connecting to database with DSN:", dsn)

	// Создаем соединение с базой данных
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Проверяем подключение
	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}
	
	log.Println("Successfully connected to database!")
	return conn, nil
}
