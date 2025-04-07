package handlers

import (
	"net/http"
	"order-service/services"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	CartService *services.CartService
}

func NewCartHandler(cartService *services.CartService) *CartHandler {
	return &CartHandler{CartService: cartService}
}

// AddToCartRequest представляет структуру запроса для добавления товара в корзину
type AddToCartRequest struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// AddToCart добавляет товар в корзину
func (h *CartHandler) AddToCart(c *gin.Context) {
	userID := c.Param("userID")

	// Получаем данные из тела запроса
	var request AddToCartRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Добавляем товар в корзину
	err := h.CartService.AddToCart(userID, request.ProductID, request.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added to cart"})
}
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userID := c.Param("userID")
	productID := c.Param("productID")

	// Удаляем товар из корзины
	err := h.CartService.RemoveFromCart(userID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product removed from cart"})
}
func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.Param("userID")

	// Получаем корзину
	cart, err := h.CartService.GetCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart": cart})
}
func (h *CartHandler) CheckoutCart(c *gin.Context) {
	userID := c.Param("userID")

	// Оформляем заказ
	order, err := h.CartService.CheckoutCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}
