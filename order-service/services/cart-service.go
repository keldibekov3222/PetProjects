package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"order-service/models"
	"order-service/repositories"
	"time"
)

type CartService struct {
	RedisClient *redis.Client
	ProductRepo *repositories.ProductRepository // Добавляем репозиторий продуктов
	OrderRepo   *repositories.OrderRepository
}

func NewCartService(redisClient *redis.Client, productRepo *repositories.ProductRepository, orderRepo *repositories.OrderRepository) *CartService {
	return &CartService{
		RedisClient: redisClient,
		ProductRepo: productRepo,
		OrderRepo:   orderRepo,
	}
}

func (s *CartService) AddToCart(userID, productID string, quantity int) error {
	ctx := context.Background()
	key := fmt.Sprintf("cart:%s", userID)

	// Проверяем, есть ли уже этот товар в корзине
	existingQuantity, err := s.RedisClient.HGet(ctx, key, productID).Int()
	if err != nil && err != redis.Nil {
		return err
	}

	// Обновляем количество товара
	totalQuantity := existingQuantity + quantity
	return s.RedisClient.HSet(ctx, key, productID, totalQuantity).Err()
}

func (s *CartService) RemoveFromCart(userID, productID string) error {
	ctx := context.Background()
	key := fmt.Sprintf("cart:%s", userID)
	return s.RedisClient.HDel(ctx, key, productID).Err()
}

func (s *CartService) GetCart(userID string) (map[string]int, error) {
	ctx := context.Background()
	key := fmt.Sprintf("cart:%s", userID)

	cartItems, err := s.RedisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	cart := make(map[string]int)
	for productID, qty := range cartItems {
		var quantity int
		json.Unmarshal([]byte(qty), &quantity)
		cart[productID] = quantity
	}
	return cart, nil
}

func (s *CartService) ClearCart(userID string) error {
	ctx := context.Background()
	key := fmt.Sprintf("cart:%s", userID)
	return s.RedisClient.Del(ctx, key).Err()
}

func generateOrderID() string {
	return uuid.New().String() // Генерируем новый UUID и преобразуем его в строку
}

func (s *CartService) CheckoutCart(userID string) (*models.Order, error) {
	// Получаем корзину из Redis
	cart, err := s.GetCart(userID)
	if err != nil {
		return nil, err
	}

	// Преобразуем cart (map[string]int) в []models.CartItem
	cartItems := []models.CartItem{}
	for productID, quantity := range cart {
		price, err := s.ProductRepo.GetProductPrice(productID)
		if err != nil {
			return nil, err
		}
		cartItems = append(cartItems, models.CartItem{
			ProductID: productID,
			Quantity:  quantity,
			Price:     price,
		})
	}

	// Создаем заказ
	order := models.Order{
		ID:         generateOrderID(), // Генерация уникального ID
		UserID:     userID,
		Items:      cartItems,
		TotalPrice: s.calculateTotalPrice(cartItems), // Метод для подсчета суммы
		Status:     "pending",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Сохраняем заказ в БД
	err = s.OrderRepo.CreateOrder(&order)
	if err != nil {
		return nil, err
	}

	// Очищаем корзину после оформления заказа
	err = s.ClearCart(userID)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
func (s *CartService) calculateTotalPrice(cartItems []models.CartItem) float64 {
	var totalPrice float64
	for _, item := range cartItems {
		// Умножаем цену товара на его количество и добавляем к общей стоимости
		totalPrice += item.Price * float64(item.Quantity)
	}
	return totalPrice
}
