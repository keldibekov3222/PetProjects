package repositories

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"order-service/models"
)

type UserRepository struct {
	DB *mongo.Database
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{DB: db}
}

// Создание нового пользователя
func (r *UserRepository) CreateUser(user *models.User) error {
	user.ID = primitive.NewObjectID()
	_, err := r.DB.Collection("users").InsertOne(context.Background(), user)
	return err
}

// Получение пользователя по email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Collection("users").FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Нет документа с таким email
		}
		return nil, err // Ошибка при запросе
	}
	return &user, nil
}
func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	// Преобразуем строковый id в ObjectId
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %v", err)
	}

	// Ищем пользователя по ObjectId
	err = r.DB.Collection("users").FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	cursor, err := r.DB.Collection("users").Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
func (r *UserRepository) UpdateUser(user *models.User) error {
	_, err := r.DB.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": user.ID}, // Идентификатор пользователя
		bson.M{
			"$set": bson.M{
				"username": user.Username,
				"email":    user.Email,
			},
		},
	)
	return err
}
func (r *UserRepository) DeleteUser(id string) error {
	// В случае MongoDB нужно использовать ObjectID, если ID хранится как строка
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.DB.Collection("users").DeleteOne(context.Background(), bson.M{"_id": objID})
	return err
}
