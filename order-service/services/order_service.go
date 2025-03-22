package services

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"order-service/models"
	"order-service/repositories"
	"time"
)

type OrderService struct {
	Repo *repositories.OrderRepository
}

func NewOrderService(repo *repositories.OrderRepository) *OrderService {
	if repo == nil {
		log.Fatal("NewOrderService: received nil repository")
	}
	return &OrderService{Repo: repo}
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
	return s.Repo.GetOrderById(id)
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
