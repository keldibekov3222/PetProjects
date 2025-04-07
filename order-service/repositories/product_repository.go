package repositories

import (
	"context"
	"fmt"
	"order-service/models"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository struct {
	db *mongo.Database
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(product *models.Product) error {
	// Генерируем новый ObjectID
	product.ID = primitive.NewObjectID()

	// Если IDString не задан, используем ObjectID в виде строки
	if product.IDString == "" {
		product.IDString = product.ID.Hex()
	}

	_, err := r.db.Collection("products").InsertOne(context.Background(), product)
	return err
}

func (r *ProductRepository) GetProductById(id string) (*models.Product, error) {
	var product models.Product

	// Проверяем, является ли ID UUID (содержит дефисы)
	if strings.Contains(id, "-") {
		// Если это UUID, ищем продукт по IDString
		err := r.db.Collection("products").FindOne(context.Background(), bson.M{"idString": id}).Decode(&product)
		if err != nil {
			return nil, fmt.Errorf("product not found: %v", err)
		}
		return &product, nil
	}

	// Если это не UUID, пробуем преобразовать в ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID format: %v", err)
	}

	err = r.db.Collection("products").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		return nil, fmt.Errorf("product not found: %v", err)
	}
	return &product, nil
}

func (r *ProductRepository) GetAllProducts() ([]models.ProductResponse, error) {
	var products []models.ProductResponse

	cursor, err := r.db.Collection("products").Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			// Пропускаем документы с некорректным _id
			continue
		}
		products = append(products, models.ProductResponse{
			ID:          product.IDString, // Используем IDString вместо ID.Hex()
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) UpdateProduct(id primitive.ObjectID, updatedProduct *models.Product) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"name":        updatedProduct.Name,
			"description": updatedProduct.Description,
			"price":       updatedProduct.Price,
			"stock":       updatedProduct.Stock,
		},
	}

	_, err := r.db.Collection("products").UpdateOne(context.Background(), filter, update)
	return err
}

func (r *ProductRepository) DeleteProduct(id primitive.ObjectID) error {
	_, err := r.db.Collection("products").DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (repo *ProductRepository) GetProductPrice(productID string) (float64, error) {
	// Проверяем, является ли ID UUID (содержит дефисы)
	if strings.Contains(productID, "-") {
		// Если это UUID, ищем продукт по IDString
		var product models.Product
		err := repo.db.Collection("products").FindOne(context.Background(), bson.M{"idString": productID}).Decode(&product)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return 0, fmt.Errorf("product not found")
			}
			return 0, err
		}
		return product.Price, nil
	}

	// Если это не UUID, пробуем преобразовать в ObjectID
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return 0, fmt.Errorf("invalid product ID format: %v", err)
	}

	var product models.Product
	err = repo.db.Collection("products").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, fmt.Errorf("product not found")
		}
		return 0, err
	}

	return product.Price, nil
}
