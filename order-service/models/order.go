package models

import "time"

type CartItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Order struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"` // UUID, внешний ключ к таблице users
	Items      []CartItem `json:"items"`
	TotalPrice float64    `json:"total_price"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
