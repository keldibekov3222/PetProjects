package repositories

import (
	"context"
	"fmt"
	"order-service/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	DB *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *UserRepository {
	return &UserRepository{DB: db}
}

// Создание нового пользователя
func (r *UserRepository) CreateUser(user *models.User) error {
	user.ID = uuid.New().String()
	query := `
		INSERT INTO users (id, username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	now := time.Now()
	_, err := r.DB.Exec(context.Background(), query,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		now,
		now,
	)
	return err
}

// Получение пользователя по email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := r.DB.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Получение пользователя по ID
func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := r.DB.QueryRow(context.Background(), query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Получение всех пользователей
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
	`
	rows, err := r.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Обновление пользователя
func (r *UserRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.DB.Exec(context.Background(), query,
		user.Username,
		user.Email,
		time.Now(),
		user.ID,
	)
	return err
}

// Удаление пользователя
func (r *UserRepository) DeleteUser(id string) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	result, err := r.DB.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
