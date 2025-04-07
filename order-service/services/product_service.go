package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"order-service/models"
	"order-service/repositories"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService struct {
	Repo        *repositories.ProductRepository
	RedisClient *redis.Client
}

func NewProductService(repo *repositories.ProductRepository, redisClient *redis.Client) *ProductService {
	return &ProductService{
		Repo:        repo,
		RedisClient: redisClient,
	}
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	return s.Repo.CreateProduct(product)
}

func (s *ProductService) GetProductById(id string) (*models.Product, error) {
	// Проверяем кэш
	cacheKey := fmt.Sprintf("product:%s", id)
	if cached, err := s.RedisClient.Get(context.Background(), cacheKey).Result(); err == nil {
		var product models.Product
		if err := json.Unmarshal([]byte(cached), &product); err == nil {
			return &product, nil
		}
	}

	// Если нет в кэше, получаем из БД
	product, err := s.Repo.GetProductById(id)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	if productJSON, err := json.Marshal(product); err == nil {
		s.RedisClient.Set(context.Background(), cacheKey, productJSON, 24*time.Hour)
	}

	return product, nil
}

func (s *ProductService) GetAllProducts() ([]models.ProductResponse, error) {
	// Проверяем кэш
	cacheKey := "products:all"
	if cached, err := s.RedisClient.Get(context.Background(), cacheKey).Result(); err == nil {
		var products []models.ProductResponse
		if err := json.Unmarshal([]byte(cached), &products); err == nil {
			return products, nil
		}
	}

	// Если нет в кэше, получаем из БД
	products, err := s.Repo.GetAllProducts()
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	if productsJSON, err := json.Marshal(products); err == nil {
		s.RedisClient.Set(context.Background(), cacheKey, productsJSON, 1*time.Hour)
	}

	return products, nil
}

func (s *ProductService) UpdateProduct(id string, updatedProduct *models.Product) (*models.Product, error) {
	// Проверяем, является ли ID UUID (содержит дефисы)
	if strings.Contains(id, "-") {
		// Если это UUID, получаем продукт по IDString
		product, err := s.Repo.GetProductById(id)
		if err != nil {
			return nil, err
		}

		// Обновляем в БД
		err = s.Repo.UpdateProduct(product.ID, updatedProduct)
		if err != nil {
			return nil, err
		}

		// Инвалидируем кэш
		cacheKey := fmt.Sprintf("product:%s", id)
		s.RedisClient.Del(context.Background(), cacheKey)
		s.RedisClient.Del(context.Background(), "products:all")

		return updatedProduct, nil
	}

	// Если это не UUID, пробуем преобразовать в ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid product ID format")
	}

	// Обновляем в БД
	err = s.Repo.UpdateProduct(objID, updatedProduct)
	if err != nil {
		return nil, err
	}

	// Инвалидируем кэш
	cacheKey := fmt.Sprintf("product:%s", id)
	s.RedisClient.Del(context.Background(), cacheKey)
	s.RedisClient.Del(context.Background(), "products:all")

	return updatedProduct, nil
}

func (s *ProductService) DeleteProduct(id string) error {
	// Проверяем, является ли ID UUID (содержит дефисы)
	if strings.Contains(id, "-") {
		// Если это UUID, получаем продукт по IDString
		product, err := s.Repo.GetProductById(id)
		if err != nil {
			return err
		}

		// Удаляем из БД
		return s.Repo.DeleteProduct(product.ID)
	}

	// Если это не UUID, пробуем преобразовать в ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid product ID format")
	}

	return s.Repo.DeleteProduct(objID)
}
