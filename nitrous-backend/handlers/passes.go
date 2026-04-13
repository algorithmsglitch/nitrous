package handlers

import (
	"net/http"
	"nitrous-backend/database"
	"nitrous-backend/middleware"

	"github.com/gin-gonic/gin"
)

type Pass struct {
	ID         string   `json:"id" db:"id"`
	Tier       string   `json:"tier" db:"tier"`
	Event      string   `json:"event" db:"event"`
	Location   string   `json:"location" db:"location"`
	Date       string   `json:"date" db:"date"`
	Category   string   `json:"category" db:"category"`
	Price      float64  `json:"price" db:"price"`
	Perks      []string `json:"perks" db:"-"`
	SpotsLeft  int      `json:"spotsLeft" db:"spots_left"`
	TotalSpots int      `json:"totalSpots" db:"total_spots"`
	Badge      *string  `json:"badge" db:"badge"`
	TierColor  string   `json:"tierColor" db:"tier_color"`
}

func PurchasePass(c *gin.Context) {
	passID := c.Param("id")
	claims := c.MustGet("claims").(*middleware.Claims)
	userID := claims.UserID

	db := database.GetDB()

	// Begin transaction — check spots and insert purchase atomically
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer tx.Rollback()

	// Lock the row and check spots remaining
	var spotsLeft int
	err = tx.QueryRow(
		`SELECT spots_left FROM passes WHERE id = $1 FOR UPDATE`,
		passID,
	).Scan(&spotsLeft)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pass not found"})
		return
	}

	if spotsLeft <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No spots remaining"})
		return
	}

	// Check user hasn't already purchased this pass
	var existing int
	err = tx.QueryRow(
		`SELECT COUNT(*) FROM pass_purchases WHERE user_id = $1 AND pass_id = $2`,
		userID, passID,
	).Scan(&existing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if existing > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You already have this pass"})
		return
	}

	// Decrement spots
	_, err = tx.Exec(
		`UPDATE passes SET spots_left = spots_left - 1 WHERE id = $1`,
		passID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reserve spot"})
		return
	}

	// Record the purchase
	_, err = tx.Exec(
		`INSERT INTO pass_purchases (user_id, pass_id) VALUES ($1, $2)`,
		userID, passID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record purchase"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete purchase"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pass purchased successfully",
		"passId":  passID,
	})
}
