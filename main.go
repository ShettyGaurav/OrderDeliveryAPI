package main

import (
	"log"

	"order-delivery-api/database"
	"order-delivery-api/handlers"
	"order-delivery-api/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitializeDatabase()
	log.Println("Database initialized successfully")

	router := gin.Default()

	router.Use(middleware.EnableCrossOrigin())
	router.Use(middleware.RequestLogger())

	api := router.Group("/api")
	{
		api.GET("/orders", handlers.FetchAllOrders)
		api.GET("/orders/:id", handlers.FetchOrderById)
		api.POST("/orders", handlers.AddNewOrder)
		api.PUT("/orders/:id", handlers.ModifyOrder)
		api.DELETE("/orders/:id", handlers.RemoveOrder)
		api.PATCH("/orders/:id/status", handlers.ChangeOrderStatus)
		api.GET("/orders/:id/track", handlers.GetOrderTracking)

		api.GET("/stats", handlers.FetchOrderStatistics)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("Starting Order Delivery API server on :9090")
	if err := router.Run(":9090"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
