package repositories

import (
	"context"
	"fmt"
	"log"
	"order-service/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type OrderRepository struct {
	DB *pgx.Conn
}

func NewOrderRepository(db *pgx.Conn) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) CreateOrder(order *models.Order) error {
	currentTime := time.Now()
	query := `INSERT INTO orders (id, user_id, total_price, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.DB.Exec(context.Background(), query, order.ID, order.UserID, order.TotalPrice, order.Status, currentTime, currentTime)
	if err != nil {
		log.Printf("error inserting order: %v", err)
		return err
	}
	return nil
}

func (r *OrderRepository) UpdateOrder(id string, updatedOrder *models.Order) (*models.Order, error) {
	// Проверяем, является ли ID валидным UUID
	orderID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("invalid order ID: %v", err)
		return nil, err
	}

	updatedOrder.UpdatedAt = time.Now()

	query := `
		UPDATE orders 
		SET user_id = $1, total_price = $2, status = $3, updated_at = $4 
		WHERE id = $5 
		RETURNING id, user_id, total_price, status, created_at, updated_at`

	// Создаём структуру для хранения обновленных данных
	var newOrder models.Order

	err = r.DB.QueryRow(context.Background(), query,
		updatedOrder.UserID, updatedOrder.TotalPrice, updatedOrder.Status, updatedOrder.UpdatedAt, orderID).
		Scan(&newOrder.ID, &newOrder.UserID, &newOrder.TotalPrice, &newOrder.Status, &newOrder.CreatedAt, &newOrder.UpdatedAt)

	if err != nil {
		log.Printf("error updating order: %v", err)
		return nil, err
	}

	return &newOrder, nil
}
func (r *OrderRepository) GetOrderById(id string) (*models.Order, error) {
	var order models.Order
	query := `
SELECT id, user_id, total_price, status, created_at, updated_at From orders WHERE id = $1`
	err := r.DB.QueryRow(context.Background(), query, id).Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		log.Printf("error getting order: %v", err)
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	rows, err := r.DB.Query(context.Background(), "SELECT id, user_id, total_price, status, created_at, updated_at FROM orders")
	if err != nil {
		fmt.Printf("error getting all orders: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status, &order.CreatedAt, &order.UpdatedAt); err != nil {
			log.Printf("error scanning order: %v", err)
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		log.Printf("error iterating over rows: %v", err)
	}
	return orders, nil
}
func (r *OrderRepository) DeleteOrder(id string) error {
	_, err := r.DB.Exec(context.Background(), "DELETE FROM orders WHERE id = $1", id)
	if err != nil {
		log.Printf("error deleting order: %v", err)
		return err
	}
	return nil
}

func (r *OrderRepository) GetOrdersByUserID(userID string) ([]models.Order, error) {
	ctx := context.Background()

	// Создаем SQL запрос для получения заказов пользователя
	query := `
		SELECT id, user_id, total_price, status, created_at, updated_at 
		FROM orders 
		WHERE user_id = $1`

	// Выполняем запрос
	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Создаем слайс для хранения результатов
	var orders []models.Order

	// Получаем все заказы
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.TotalPrice,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
