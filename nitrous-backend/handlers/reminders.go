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
// SetReminder creates a reminder for the authenticated user.
func SetReminder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.SetReminderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.RemindAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reminder time must be in the future"})
		return
	}

	eventExists := false
	for _, event := range database.Events {
		if event.ID == req.EventID {
=======
// SetReminder creates a reminder for an event (auth required)
func SetReminder(c *gin.Context) {
	eventID := c.Param("id")
	userID := c.GetString("userID")

	// Check event exists
	eventExists := false
	for _, e := range database.Events {
		if e.ID == eventID {
>>>>>>> Stashed changes
			eventExists = true
			break
		}
	}
<<<<<<< Updated upstream

=======
>>>>>>> Stashed changes
	if !eventExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

<<<<<<< Updated upstream
	reminder := models.Reminder{
		ID:        uuid.New().String(),
		UserID:    userID.(string),
		EventID:   req.EventID,
		Message:   req.Message,
		RemindAt:  req.RemindAt,
=======
	// Check not already set
	for _, r := range database.Reminders {
		if r.UserID == userID && r.EventID == eventID {
			c.JSON(http.StatusConflict, gin.H{"error": "Reminder already set for this event"})
			return
		}
	}

	reminder := models.Reminder{
		ID:        uuid.New().String(),
		UserID:    userID,
		EventID:   eventID,
>>>>>>> Stashed changes
		CreatedAt: time.Now(),
	}

	database.Reminders = append(database.Reminders, reminder)

<<<<<<< Updated upstream
	c.JSON(http.StatusCreated, reminder)
}

// DeleteReminder deletes one reminder owned by the authenticated user.
func DeleteReminder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	reminderID := c.Param("id")

	for i, reminder := range database.Reminders {
		if reminder.ID == reminderID {
			if reminder.UserID != userID.(string) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
				return
			}

			database.Reminders = append(database.Reminders[:i], database.Reminders[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Reminder deleted"})
=======
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Reminder set",
		"reminder": reminder,
	})
}

// DeleteReminder removes a reminder
func DeleteReminder(c *gin.Context) {
	eventID := c.Param("id")
	userID := c.GetString("userID")

	for i, r := range database.Reminders {
		if r.UserID == userID && r.EventID == eventID {
			database.Reminders = append(database.Reminders[:i], database.Reminders[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Reminder removed"})
>>>>>>> Stashed changes
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Reminder not found"})
}

<<<<<<< Updated upstream
// GetMyReminders returns all reminders for the authenticated user.
func GetMyReminders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var reminders []models.Reminder
	for _, reminder := range database.Reminders {
		if reminder.UserID == userID.(string) {
			reminders = append(reminders, reminder)
=======
// GetMyReminders returns all reminders for the authenticated user
func GetMyReminders(c *gin.Context) {
	userID := c.GetString("userID")

	var userReminders []models.Reminder
	for _, r := range database.Reminders {
		if r.UserID == userID {
			userReminders = append(userReminders, r)
>>>>>>> Stashed changes
		}
	}

	c.JSON(http.StatusOK, gin.H{
<<<<<<< Updated upstream
		"reminders": reminders,
		"count":     len(reminders),
=======
		"reminders": userReminders,
		"count":     len(userReminders),
>>>>>>> Stashed changes
	})
}
