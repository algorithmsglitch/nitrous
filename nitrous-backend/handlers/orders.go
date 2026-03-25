package handlers

import (
	"net/http"
	"nitrous-backend/database"
	"nitrous-backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

<<<<<<< Updated upstream
// CreateOrder creates a merch order for the authenticated user.
func CreateOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
=======
// CreateOrder places a merch order (auth required)
func CreateOrder(c *gin.Context) {
	userID := c.GetString("userID")
>>>>>>> Stashed changes

	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

<<<<<<< Updated upstream
	var merchItem *models.MerchItem
	for _, item := range database.MerchItems {
		if item.ID == req.MerchItemID {
			itemCopy := item
			merchItem = &itemCopy
			break
		}
	}

	if merchItem == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merch item not found"})
		return
	}

	order := models.Order{
		ID:          uuid.New().String(),
		UserID:      userID.(string),
		MerchItemID: req.MerchItemID,
		Quantity:    req.Quantity,
		UnitPrice:   merchItem.Price,
		TotalPrice:  merchItem.Price * float64(req.Quantity),
		Status:      "created",
		CreatedAt:   time.Now(),
=======
	// Validate each item exists and calculate total
	var total float64
	for i, item := range req.Items {
		found := false
		for _, m := range database.MerchItems {
			if m.ID == item.MerchID {
				found = true
				req.Items[i].Name = m.Name
				req.Items[i].Price = m.Price
				total += m.Price * float64(item.Quantity)
				break
			}
		}
		if !found {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Merch item not found: " + item.MerchID})
			return
		}
	}

	order := models.Order{
		ID:        uuid.New().String(),
		UserID:    userID,
		Items:     req.Items,
		Total:     total,
		Status:    "pending",
		CreatedAt: time.Now(),
>>>>>>> Stashed changes
	}

	database.Orders = append(database.Orders, order)

<<<<<<< Updated upstream
	c.JSON(http.StatusCreated, order)
}

// GetMyOrders returns all orders for the authenticated user.
func GetMyOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var orders []models.Order
	for _, order := range database.Orders {
		if order.UserID == userID.(string) {
			orders = append(orders, order)
=======
	c.JSON(http.StatusCreated, gin.H{
		"message": "Order placed successfully",
		"order":   order,
	})
}

// GetMyOrders returns all orders for the authenticated user
func GetMyOrders(c *gin.Context) {
	userID := c.GetString("userID")

	var userOrders []models.Order
	for _, o := range database.Orders {
		if o.UserID == userID {
			userOrders = append(userOrders, o)
>>>>>>> Stashed changes
		}
	}

	c.JSON(http.StatusOK, gin.H{
<<<<<<< Updated upstream
		"orders": orders,
		"count":  len(orders),
	})
}

// GetOrderByID returns one order if it belongs to the authenticated user.
func GetOrderByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	orderID := c.Param("id")

	for _, order := range database.Orders {
		if order.ID == orderID {
			if order.UserID != userID.(string) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
				return
			}

			c.JSON(http.StatusOK, order)
=======
		"orders": userOrders,
		"count":  len(userOrders),
	})
}

// GetOrderByID returns a single order (must belong to authenticated user)
func GetOrderByID(c *gin.Context) {
	orderID := c.Param("id")
	userID := c.GetString("userID")

	for _, o := range database.Orders {
		if o.ID == orderID {
			if o.UserID != userID {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
				return
			}
			c.JSON(http.StatusOK, o)
>>>>>>> Stashed changes
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
}
