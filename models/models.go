package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	StatusPending        OrderStatus = "pending"
	StatusConfirmed      OrderStatus = "confirmed"
	StatusPreparing      OrderStatus = "preparing"
	StatusOutForDelivery OrderStatus = "out_for_delivery"
	StatusDelivered      OrderStatus = "delivered"
	StatusCancelled      OrderStatus = "cancelled"
)

type OrderItem struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	OrderID  uint    `json:"order_id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type Order struct {
	gorm.Model
	CustomerName      string      `json:"customer_name"`
	CustomerAddress   string      `json:"customer_address"`
	CustomerPhone     string      `json:"customer_phone"`
	CustomerEmail     string      `json:"customer_email"`
	Items             []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	TotalAmount       float64     `json:"total_amount"`
	Status            OrderStatus `json:"status" gorm:"default:pending"`
	DeliveryPerson    string      `json:"delivery_person"`
	DeliveryPhone     string      `json:"delivery_phone"`
	Notes             string      `json:"notes"`
	EstimatedDelivery time.Time   `json:"estimated_delivery"`
}

type TrackingInfo struct {
	OrderID           uint        `json:"order_id"`
	Status            OrderStatus `json:"status"`
	DeliveryPerson    string      `json:"delivery_person"`
	DeliveryPhone     string      `json:"delivery_phone"`
	EstimatedDelivery time.Time   `json:"estimated_delivery"`
	UpdatedAt         time.Time   `json:"updated_at"`
	StatusHistory     []StatusLog `json:"status_history"`
}

type StatusLog struct {
	Status    OrderStatus `json:"status"`
	Timestamp time.Time   `json:"timestamp"`
	Message   string      `json:"message"`
}

type CreateOrderRequest struct {
	CustomerName      string      `json:"customer_name" binding:"required"`
	CustomerAddress   string      `json:"customer_address" binding:"required"`
	CustomerPhone     string      `json:"customer_phone" binding:"required"`
	CustomerEmail     string      `json:"customer_email"`
	Items             []OrderItem `json:"items" binding:"required"`
	Notes             string      `json:"notes"`
	EstimatedDelivery time.Time   `json:"estimated_delivery"`
}

type UpdateOrderRequest struct {
	CustomerName      string      `json:"customer_name"`
	CustomerAddress   string      `json:"customer_address"`
	CustomerPhone     string      `json:"customer_phone"`
	CustomerEmail     string      `json:"customer_email"`
	Items             []OrderItem `json:"items"`
	Notes             string      `json:"notes"`
	DeliveryPerson    string      `json:"delivery_person"`
	DeliveryPhone     string      `json:"delivery_phone"`
	EstimatedDelivery time.Time   `json:"estimated_delivery"`
}

type UpdateStatusRequest struct {
	Status  OrderStatus `json:"status" binding:"required"`
	Message string      `json:"message"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type OrderStats struct {
	TotalOrders     int64   `json:"total_orders"`
	PendingOrders   int64   `json:"pending_orders"`
	DeliveredOrders int64   `json:"delivered_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
}
