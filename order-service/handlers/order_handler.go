package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"order-service/models"
	"order-service/services"
)

type OrderHandler struct {
	Service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{Service: service}
}

func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	var request struct {
		UserID     string  `json:"user_id"`
		TotalPrice float64 `json:"total_price"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil { // Используем ShouldBindJSON вместо ShouldBind
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.Service.CreateOrder(request.UserID, request.TotalPrice)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrderById(ctx *gin.Context) {
	id := ctx.Param("id")

	order, err := h.Service.GetOrderById(id) // Исправлено на `h.Service`
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Order with id %s not found", id)})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetAllOrders(ctx *gin.Context) {
	orders, err := h.Service.GetAllOrders() // Исправлено на `h.Service`
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) UpdateOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	var updatedOrder models.Order

	if err := ctx.ShouldBindJSON(&updatedOrder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Правильное присваивание с указателем
	updatedOrderPtr, err := h.Service.UpdateOrder(id, &updatedOrder)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	// Разыменовываем указатель перед отправкой в JSON
	ctx.JSON(http.StatusOK, *updatedOrderPtr)
}

func (h *OrderHandler) DeleteOrder(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.Service.DeleteOrder(id) // Исправлено на `h.Service`
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Order with id %s not found", id)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
