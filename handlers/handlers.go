package handlers

import (
	"net/http"
	"sort"
	"strconv"

	"order-delivery-api/database"
	"order-delivery-api/models"

	"github.com/gin-gonic/gin"
)

func FetchAllOrders(c *gin.Context) {
	orders := database.FetchAllOrders()

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].CreatedAt.After(orders[j].CreatedAt)
	})

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Orders retrieved successfully",
		Data:    orders,
	})
}

func FetchOrderById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid order ID",
		})
		return
	}

	order, exists := database.FetchOrderById(uint(id))
	if !exists {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Order not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Order retrieved successfully",
		Data:    order,
	})
}

func AddNewOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	order := database.AddNewOrder(req)

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Order created successfully",
		Data:    order,
	})
}

func ModifyOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid order ID",
		})
		return
	}

	var req models.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	order, exists := database.ModifyOrder(uint(id), req)
	if !exists {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Order not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Order updated successfully",
		Data:    order,
	})
}

func RemoveOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid order ID",
		})
		return
	}

	deleted := database.RemoveOrder(uint(id))
	if !deleted {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Order not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Order deleted successfully",
	})
}

func ChangeOrderStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid order ID",
		})
		return
	}

	var req models.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	validStatuses := map[models.OrderStatus]bool{
		models.StatusPending:        true,
		models.StatusConfirmed:      true,
		models.StatusPreparing:      true,
		models.StatusOutForDelivery: true,
		models.StatusDelivered:      true,
		models.StatusCancelled:      true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid status",
		})
		return
	}

	order, exists := database.ChangeOrderStatus(uint(id), req.Status)
	if !exists {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Order not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Order status updated successfully",
		Data:    order,
	})
}

func GetOrderTracking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid order ID",
		})
		return
	}

	order, exists := database.FetchOrderById(uint(id))
	if !exists {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Order not found",
		})
		return
	}

	statusHistory := []models.StatusLog{
		{Status: models.StatusPending, Timestamp: order.CreatedAt, Message: "Order placed"},
	}

	statusOrder := []models.OrderStatus{
		models.StatusConfirmed,
		models.StatusPreparing,
		models.StatusOutForDelivery,
		models.StatusDelivered,
	}

	currentIndex := -1
	for i, s := range statusOrder {
		if order.Status == s {
			currentIndex = i
			break
		}
	}

	for i := 0; i <= currentIndex && i < len(statusOrder); i++ {
		statusHistory = append(statusHistory, models.StatusLog{
			Status:    statusOrder[i],
			Timestamp: order.UpdatedAt,
			Message:   getMessageForStatus(statusOrder[i]),
		})
	}

	trackingInfo := models.TrackingInfo{
		OrderID:           order.ID,
		Status:            order.Status,
		DeliveryPerson:    order.DeliveryPerson,
		DeliveryPhone:     order.DeliveryPhone,
		EstimatedDelivery: order.EstimatedDelivery,
		UpdatedAt:         order.UpdatedAt,
		StatusHistory:     statusHistory,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Tracking information retrieved",
		Data:    trackingInfo,
	})
}

func FetchOrderStatistics(c *gin.Context) {
	stats := database.FetchStatistics()

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Statistics retrieved successfully",
		Data:    stats,
	})
}

func getMessageForStatus(status models.OrderStatus) string {
	messages := map[models.OrderStatus]string{
		models.StatusPending:        "Order placed",
		models.StatusConfirmed:      "Order confirmed by restaurant",
		models.StatusPreparing:      "Order is being prepared",
		models.StatusOutForDelivery: "Order is out for delivery",
		models.StatusDelivered:      "Order delivered successfully",
		models.StatusCancelled:      "Order cancelled",
	}
	return messages[status]
}
