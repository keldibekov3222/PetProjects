package services

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"order-service/models"
	"order-service/repositories"
)

type ProductService struct {
	Repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{
		Repo: repo,
	}
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	return s.Repo.CreateProduct(product)
}

func (s *ProductService) GetProductById(id string) (*models.Product, error) {
	return s.Repo.GetProductById(id)
}
func (s *ProductService) GetAllProducts() ([]models.ProductResponse, error) {
	return s.Repo.GetAllProducts()
}
func (s *ProductService) UpdateProduct(id string, updatedProduct *models.Product) (*models.Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid product ID format")
	}

	// Обновление без изменения поля UpdatedAt
	err = s.Repo.UpdateProduct(objID, updatedProduct)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

// ❌ Удаление продукта
func (s *ProductService) DeleteProduct(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid product ID format")
	}

	return s.Repo.DeleteProduct(objID)
}
