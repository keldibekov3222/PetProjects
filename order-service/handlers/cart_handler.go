package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"order-service/services"
	"strconv"
)

type CartHandler struct {
	CartService *services.CartService
}

func NewCartHandler(cartService *services.CartService) *CartHandler {
	return &CartHandler{CartService: cartService}
}

// AddToCart добавляет товар в корзину
func (h *CartHandler) AddToCart(c *gin.Context) {
	userID := c.Param("userID")
	productID := c.Param("productID")
	quantity := c.DefaultQuery("quantity", "1") // Количество товара, по умолчанию 1

	// Преобразуем количество из строки в число
	qty, err := strconv.Atoi(quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
		return
	}

	// Добавляем товар в корзину
	err = h.CartService.AddToCart(userID, productID, qty)
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
