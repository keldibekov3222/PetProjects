package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"order-service/models"
	"order-service/repositories"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type OrderService struct {
	Repo        *repositories.OrderRepository
	RedisClient *redis.Client
}

// OrderStats представляет статистику заказов
type OrderStats struct {
	TotalOrders    int64   `json:"total_orders"`
	TotalRevenue   float64 `json:"total_revenue"`
	AverageOrder   float64 `json:"average_order"`
	OrdersPerMonth int64   `json:"orders_per_month"`
}

func NewOrderService(repo *repositories.OrderRepository, redisClient *redis.Client) *OrderService {
	if repo == nil {
		log.Fatal("NewOrderService: received nil repository")
	}
	return &OrderService{
		Repo:        repo,
		RedisClient: redisClient,
	}
}

func (s *OrderService) CreateOrder(userID string, totalPrice float64) (*models.Order, error) {

	if s.Repo == nil {
		return nil, errors.New("order repository is not initialized")
	}

	order := &models.Order{
		ID:         uuid.New().String(),
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}
	err := s.Repo.CreateOrder(order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) GetOrderById(id string) (*models.Order, error) {
	// Проверяем кэш
	cacheKey := fmt.Sprintf("order:%s", id)
	if cached, err := s.RedisClient.Get(context.Background(), cacheKey).Result(); err == nil {
		var order models.Order
		if err := json.Unmarshal([]byte(cached), &order); err == nil {
			return &order, nil
		}
	}

	// Если нет в кэше, получаем из БД
	order, err := s.Repo.GetOrderById(id)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	if orderJSON, err := json.Marshal(order); err == nil {
		s.RedisClient.Set(context.Background(), cacheKey, orderJSON, 1*time.Hour)
	}

	return order, nil
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.Repo.GetAllOrders()
}

func (s *OrderService) DeleteOrder(id string) error {
	return s.Repo.DeleteOrder(id)
}

func (s *OrderService) UpdateOrder(id string, updatedOrder *models.Order) (*models.Order, error) {
	return s.Repo.UpdateOrder(id, updatedOrder)
}

// Кэширование последних заказов пользователя
func (s *OrderService) GetUserOrders(userID string) ([]models.Order, error) {
	cacheKey := fmt.Sprintf("user_orders:%s", userID)
	if cached, err := s.RedisClient.Get(context.Background(), cacheKey).Result(); err == nil {
		var orders []models.Order
		if err := json.Unmarshal([]byte(cached), &orders); err == nil {
			return orders, nil
		}
	}

	orders, err := s.Repo.GetOrdersByUserID(userID)
	if err != nil {
		return nil, err
	}

	if ordersJSON, err := json.Marshal(orders); err == nil {
		s.RedisClient.Set(context.Background(), cacheKey, ordersJSON, 30*time.Minute)
	}

	return orders, nil
}

func (s *OrderService) GetOrderStatistics() (*OrderStats, error) {
	cacheKey := "order:statistics"
	if cached, err := s.RedisClient.Get(context.Background(), cacheKey).Result(); err == nil {
		var stats OrderStats
		if err := json.Unmarshal([]byte(cached), &stats); err == nil {
			return &stats, nil
		}
	}

	// Получаем статистику из БД
	orders, err := s.Repo.GetAllOrders()
	if err != nil {
		return nil, err
	}

	// Вычисляем статистику
	stats := &OrderStats{
		TotalOrders: int64(len(orders)),
	}

	var totalRevenue float64
	for _, order := range orders {
		totalRevenue += order.TotalPrice
	}
	stats.TotalRevenue = totalRevenue

	if stats.TotalOrders > 0 {
		stats.AverageOrder = totalRevenue / float64(stats.TotalOrders)
	}

	// Сохраняем в кэш
	if statsJSON, err := json.Marshal(stats); err == nil {
		s.RedisClient.Set(context.Background(), cacheKey, statsJSON, 1*time.Hour)
	}

	return stats, nil
}
