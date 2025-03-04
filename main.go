package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// โครงสร้างข้อมูล
type Coffee struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

type Order struct {
	OrderID           string `json:"order_id"`
	CoffeeID          string `json:"coffee_id"`
	Quantity          int    `json:"quantity"`
	CreatedAt         string `json:"created_at"`
	EstimatedDelivery string `json:"estimated_delivery"`
	Status            string `json:"status"`
}

var coffees = []Coffee{
	{ID: "c001", Name: "Espresso", Type: "Espresso", Price: 60, Description: "เข้มข้น กลมกล่อม"},
	{ID: "c002", Name: "Americano", Type: "Espresso", Price: 65, Description: "เอสเพรสโซ่ผสมน้ำร้อน"},
	{ID: "c003", Name: "Latte", Type: "Latte", Price: 75, Description: "เอสเพรสโซ่ผสมนมร้อน"},
}

var orders = []Order{
	{
		OrderID:           "o12345",
		CoffeeID:          "c003",
		Quantity:          2,
		CreatedAt:         "2025-03-03T15:30:45+07:00",
		EstimatedDelivery: "2025-03-03T15:40:45+07:00",
		Status:            "In Progress",
	},
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		c.Next()
		latency := time.Since(startTime)
		log.Printf("| %s | %s | %v |", c.Request.Method, path, latency)
	}
}

func getAllCoffees(c *gin.Context) {
	c.JSON(200, coffees)
}
func getCoffeeByID(c *gin.Context) {
	id := c.Param("id")
	for _, coffee := range coffees {
		if coffee.ID == id {
			c.JSON(200, coffee)
			return
		}
	}
	c.JSON(404, gin.H{"error": "Coffee not found"})
}

func searchCoffeesHandler(c *gin.Context) {
	nameQuery := c.Query("name")
	typeQuery := c.Query("type")
	results := []Coffee{}

	for _, coffee := range coffees {
		if (nameQuery == "" || strings.Contains(strings.ToLower(coffee.Name), strings.ToLower(nameQuery))) &&
			(typeQuery == "" || strings.EqualFold(coffee.Type, typeQuery)) {
			results = append(results, coffee)
		}
	}

	c.JSON(http.StatusOK, gin.H{"coffees": results})
}

func findCoffeeByID(id string) *Coffee {
	for _, coffee := range coffees {
		if coffee.ID == id {
			return &coffee
		}
	}
	return nil
}

func createOrderHandler(c *gin.Context) {
	var newOrder struct {
		CoffeeID string `json:"coffee_id"`
		Quantity int    `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&newOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if newOrder.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity, must be greater than 0"})
		return
	}

	coffee := findCoffeeByID(newOrder.CoffeeID)
	if coffee == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coffee not found"})
		return
	}

	orderID := uuid.New().String()
	createdAt := time.Now().Format(time.RFC3339)
	estimatedDelivery := time.Now().Add(10 * time.Minute).Format(time.RFC3339)

	order := Order{
		OrderID:           orderID,
		CoffeeID:          newOrder.CoffeeID,
		Quantity:          newOrder.Quantity,
		CreatedAt:         createdAt,
		EstimatedDelivery: estimatedDelivery,
		Status:            "Pending",
	}

	orders = append(orders, order)

	c.JSON(http.StatusCreated, order)
}

func getOrderByID(c *gin.Context) {
	orderID := c.Param("order_id")

	for _, order := range orders {
		if order.OrderID == orderID {
			c.JSON(http.StatusOK, order)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
}

func main() {
	r := gin.Default()
	r.Use(LoggingMiddleware()) //MiddleWare

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "welcome!")
	})

	r.GET("/coffees", getAllCoffees)               // 1.
	r.GET("/coffees/:id", getCoffeeByID)           // 2.
	r.GET("/coffees/search", searchCoffeesHandler) // 3.
	r.POST("/orders", createOrderHandler)          // 4.
	r.GET("/orders/:order_id", getOrderByID)       // 5.

	r.Run(":8080")
}
