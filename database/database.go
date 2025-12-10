package database

import (
	"fmt"
	"log"
	"order-delivery-api/models"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func InitializeDatabase() {
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "1234")
	dbname := getEnv("DB_NAME", "order_delivery")
	port := getEnv("DB_PORT", "5432")
	sslmode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(&models.Order{}, &models.OrderItem{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("PostgreSQL database connected and migrated successfully")
}

func GetDatabaseInstance() *gorm.DB {
	return db
}

func FetchAllOrders() []models.Order {
	var orders []models.Order
	db.Preload("Items").Find(&orders)
	return orders
}

func FetchOrderById(id uint) (*models.Order, bool) {
	var order models.Order
	result := db.Preload("Items").First(&order, id)
	if result.Error != nil {
		return nil, false
	}
	return &order, true
}

func AddNewOrder(req models.CreateOrderRequest) *models.Order {
	var total float64
	for _, item := range req.Items {
		total += item.Price * float64(item.Quantity)
	}

	order := &models.Order{
		CustomerName:      req.CustomerName,
		CustomerAddress:   req.CustomerAddress,
		CustomerPhone:     req.CustomerPhone,
		CustomerEmail:     req.CustomerEmail,
		Items:             req.Items,
		TotalAmount:       total,
		Status:            models.StatusPending,
		Notes:             req.Notes,
		EstimatedDelivery: time.Now().Add(45 * time.Minute),
	}

	db.Create(order)
	return order
}

func ModifyOrder(id uint, req models.UpdateOrderRequest) (*models.Order, bool) {
	var order models.Order
	result := db.First(&order, id)
	if result.Error != nil {
		return nil, false
	}

	if req.CustomerName != "" {
		order.CustomerName = req.CustomerName
	}
	if req.CustomerAddress != "" {
		order.CustomerAddress = req.CustomerAddress
	}
	if req.CustomerPhone != "" {
		order.CustomerPhone = req.CustomerPhone
	}
	if req.CustomerEmail != "" {
		order.CustomerEmail = req.CustomerEmail
	}
	if len(req.Items) > 0 {
		db.Where("order_id = ?", id).Delete(&models.OrderItem{})
		order.Items = req.Items
		var total float64
		for _, item := range req.Items {
			total += item.Price * float64(item.Quantity)
		}
		order.TotalAmount = total
	}
	if req.Notes != "" {
		order.Notes = req.Notes
	}
	if req.DeliveryPerson != "" {
		order.DeliveryPerson = req.DeliveryPerson
	}
	if req.DeliveryPhone != "" {
		order.DeliveryPhone = req.DeliveryPhone
	}
	if !req.EstimatedDelivery.IsZero() {
		order.EstimatedDelivery = req.EstimatedDelivery
	}

	db.Save(&order)
	return &order, true
}

func RemoveOrder(id uint) bool {
	result := db.Delete(&models.Order{}, id)
	return result.RowsAffected > 0
}

func ChangeOrderStatus(id uint, status models.OrderStatus) (*models.Order, bool) {
	var order models.Order
	result := db.First(&order, id)
	if result.Error != nil {
		return nil, false
	}

	order.Status = status
	db.Save(&order)
	return &order, true
}

func FetchStatistics() models.OrderStats {
	var stats models.OrderStats

	db.Model(&models.Order{}).Count(&stats.TotalOrders)
	db.Model(&models.Order{}).Where("status = ?", models.StatusPending).Count(&stats.PendingOrders)
	db.Model(&models.Order{}).Where("status = ?", models.StatusDelivered).Count(&stats.DeliveredOrders)

	var totalRevenue float64
	db.Model(&models.Order{}).Select("COALESCE(SUM(total_amount), 0)").Scan(&totalRevenue)
	stats.TotalRevenue = totalRevenue

	return stats
}
