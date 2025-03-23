package services

import (
	"context"
	"encoding/json"
	"fmt"
	"order-service/models"
	"order-service/repositories"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo        *repositories.UserRepository
	RedisClient *redis.Client
}

func NewUserService(repo *repositories.UserRepository, redisClient *redis.Client) *UserService {
	return &UserService{
		Repo:        repo,
		RedisClient: redisClient,
	}
}

// Регистрация пользователя
func (s *UserService) Register(username, email, password string) (*models.User, error) {
	// Проверка на наличие пользователя с таким email
	existingUser, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		// Если ошибка не связана с отсутствием пользователя, вернем ее
		return nil, fmt.Errorf("error checking for existing user: %v", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	// Сохранение пользователя в базе данных
	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	err = s.Repo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("error saving user to database: %v", err)
	}

	return user, nil
}

// Вход в систему
func (s *UserService) Login(email, password string) (string, error) {
	// Получение пользователя по email
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	// Генерация JWT токена
	token, err := s.generateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Генерация JWT токена
func (s *UserService) generateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret_key"))
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	// Проверяем кэш
	cacheKey := fmt.Sprintf("user:%s", id)
	if cached, err := s.RedisClient.Get(context.Background(), cacheKey).Result(); err == nil {
		var user models.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			return &user, nil
		}
	}

	// Если нет в кэше, получаем из БД
	user, err := s.Repo.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Сохраняем в кэш
	if userJSON, err := json.Marshal(user); err == nil {
		s.RedisClient.Set(context.Background(), cacheKey, userJSON, 12*time.Hour)
	}

	return user, nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.Repo.GetAllUsers()
}

func (s *UserService) UpdateUser(id string, updatedUser *models.User) (*models.User, error) {
	// Проверка на наличие пользователя
	user, err := s.Repo.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Обновляем информацию о пользователе
	user.Username = updatedUser.Username
	user.Email = updatedUser.Email
	// Можно обновить другие поля при необходимости

	err = s.Repo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(id string) error {
	// Проверяем кэш
	cacheKey := fmt.Sprintf("user:%s", id)
	s.RedisClient.Del(context.Background(), cacheKey)

	// Удаляем пользователя
	err := s.Repo.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	return nil
}
