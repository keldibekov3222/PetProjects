package db

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"order-service/config"

	"github.com/jackc/pgx/v5"
)

// Migrate выполняет миграции базы данных
func Migrate(db *pgx.Conn) error {
	// Определяем путь к директории с миграциями
	// Сначала проверяем абсолютный путь для Docker
	migrationsDir := "/app/db/migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		// Если не найдено, используем относительный путь для локальной разработки
		migrationsDir = "db/migrations"
	}

	log.Printf("Используется директория миграций: %s", migrationsDir)

	// Проверяем существование всех необходимых таблиц
	requiredTables := []string{"users", "products", "orders", "payments"}
	var missingTables []string

	for _, table := range requiredTables {
		var exists bool
		err := db.QueryRow(context.Background(), `
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = $1
			);
		`, table).Scan(&exists)

		if err != nil {
			return fmt.Errorf("ошибка при проверке существования таблицы %s: %v", table, err)
		}

		if !exists {
			missingTables = append(missingTables, table)
		}
	}

	// Если все таблицы существуют, пропускаем миграции
	if len(missingTables) == 0 {
		log.Println("Все необходимые таблицы уже существуют, пропускаем миграции")
		return nil
	}

	log.Printf("Отсутствуют таблицы: %v", missingTables)

	// Получаем список файлов миграций
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("ошибка при чтении директории миграций: %v", err)
	}

	// Фильтруем только .sql файлы и сортируем их
	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	log.Printf("Найдено файлов миграций: %d", len(migrationFiles))

	// Применяем каждую миграцию
	for _, file := range migrationFiles {
		log.Printf("Применение миграции: %s", file)
		content, err := ioutil.ReadFile(filepath.Join(migrationsDir, file))
		if err != nil {
			return fmt.Errorf("ошибка при чтении файла миграции %s: %v", file, err)
		}

		// Выполняем SQL-запросы из файла
		_, err = db.Exec(context.Background(), string(content))
		if err != nil {
			return fmt.Errorf("ошибка при выполнении миграции %s: %v", file, err)
		}
		log.Printf("Миграция %s успешно применена", file)
	}

	log.Println("Все миграции успешно применены")
	return nil
}

func MigrateConfig(cfg *config.Config) error {
	log.Println("Starting database migration...")

	// Формируем строку подключения к базе данных
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	// Подключаемся к базе данных
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// Применяем миграции
	if err := Migrate(conn); err != nil {
		return fmt.Errorf("error during migration: %v", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}
