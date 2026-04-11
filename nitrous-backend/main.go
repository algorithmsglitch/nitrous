package main

import (
	"log"
	"nitrous-backend/config"
	"nitrous-backend/database"
	"nitrous-backend/handlers"
	"nitrous-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	config.LoadConfig()

	// Initialize database
	database.InitDB()
	defer database.CloseDB()

	database.InitNewCollections() // ← new: seeds Teams, Streams
	defer database.CloseDB()

	// Start WebSocket hub
	go handlers.RunHub()

	// Dev only: simulate live telemetry updates every 5s
	// Comment out in production and replace with real data source
	go handlers.SimulateTelemetry()

	// Create Gin router
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://nitrous.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Nitrous API is running"})
	})

	// API routes
	api := r.Group("/api")
	{
		// Events
	// WebSocket — live stream telemetry
	r.GET("/ws/streams", handlers.StreamsWS)

	// API routes
	api := r.Group("/api")
	{
		// ── Events ──────────────────────────────────────────────────────────
		events := api.Group("/events")
		{
			events.GET("", handlers.GetEvents)
			events.GET("/live", handlers.GetLiveEvents)
			events.GET("/:id", handlers.GetEventByID)
<<<<<<< HEAD
<<<<<<< Updated upstream
=======
			events.POST("/:id/remind", middleware.AuthMiddleware(), handlers.SetReminder)
			events.DELETE("/:id/remind", middleware.AuthMiddleware(), handlers.DeleteReminder)
>>>>>>> Stashed changes
			events.POST("", middleware.AuthMiddleware(), handlers.CreateEvent)
			events.PUT("/:id", middleware.AuthMiddleware(), handlers.UpdateEvent)
			events.DELETE("/:id", middleware.AuthMiddleware(), handlers.DeleteEvent)
=======
			events.POST("", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.CreateEvent)
			events.PUT("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.UpdateEvent)
			events.DELETE("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteEvent)
>>>>>>> edda093f8c6ae5b6d629683c166654cfa599f29e
		}

<<<<<<< Updated upstream
		// Categories
=======
		// ── Streams ──────────────────────────────────────────────────────────
		streams := api.Group("/streams")
		{
			streams.GET("", handlers.GetStreams)
			streams.GET("/:id", handlers.GetStreamByID)
		}

		// ── Categories ───────────────────────────────────────────────────────
>>>>>>> Stashed changes
		categories := api.Group("/categories")
		{
			categories.GET("", handlers.GetCategories)
			categories.GET("/:slug", handlers.GetCategoryBySlug)
			categories.POST("", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.CreateCategory)
			categories.PUT("/:slug", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.UpdateCategory)
			categories.DELETE("/:slug", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteCategory)
		}

<<<<<<< Updated upstream
		// Journeys
=======
		// ── Journeys ─────────────────────────────────────────────────────────
>>>>>>> Stashed changes
		journeys := api.Group("/journeys")
		{
			journeys.GET("", handlers.GetJourneys)
			journeys.GET("/:id", handlers.GetJourneyByID)
			journeys.POST("", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.CreateJourney)
			journeys.PUT("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.UpdateJourney)
			journeys.DELETE("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteJourney)
			journeys.POST("/:id/book", middleware.AuthMiddleware(), handlers.BookJourney)
		}

<<<<<<< Updated upstream
		// Merch
=======
		// ── Merch ────────────────────────────────────────────────────────────
>>>>>>> Stashed changes
		merch := api.Group("/merch")
		{
			merch.GET("", handlers.GetMerchItems)
			merch.GET("/:id", handlers.GetMerchItemByID)
		}

<<<<<<< Updated upstream
		// Teams
		teams := api.Group("/teams")
		{
			teams.GET("", handlers.GetTeams)
			teams.GET("/:id", handlers.GetTeamByID)
			teams.POST("", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.CreateTeam)
			teams.PUT("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.UpdateTeam)
			teams.DELETE("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteTeam)
			teams.POST("/:id/follow", middleware.AuthMiddleware(), handlers.FollowTeam)
			teams.POST("/:id/unfollow", middleware.AuthMiddleware(), handlers.UnfollowTeam)
		}

		// Streams
		streams := api.Group("/streams")
		{
			streams.GET("", handlers.GetStreams)
			streams.GET("/:id", handlers.GetStreamByID)
			streams.POST("", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.CreateStream)
			streams.PUT("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.UpdateStream)
			streams.DELETE("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteStream)
			streams.GET("/ws", handlers.StreamsWS)
		}

		// Reminders
		reminders := api.Group("/reminders")
		{
			reminders.GET("", middleware.AuthMiddleware(), handlers.GetMyReminders)
			reminders.POST("", middleware.AuthMiddleware(), handlers.SetReminder)
			reminders.DELETE("/:id", middleware.AuthMiddleware(), handlers.DeleteReminder)
		}

		// Orders
		orders := api.Group("/orders")
		{
			orders.GET("", middleware.AuthMiddleware(), handlers.GetMyOrders)
			orders.POST("", middleware.AuthMiddleware(), handlers.CreateOrder)
			orders.GET("/:id", middleware.AuthMiddleware(), handlers.GetOrderByID)
		}

		// Auth
=======
		// ── Orders ───────────────────────────────────────────────────────────
		orders := api.Group("/orders")
		orders.Use(middleware.AuthMiddleware())
		{
			orders.POST("", handlers.CreateOrder)
			orders.GET("", handlers.GetMyOrders)
			orders.GET("/:id", handlers.GetOrderByID)
		}

		// ── Auth ─────────────────────────────────────────────────────────────
>>>>>>> Stashed changes
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.GET("/me", middleware.AuthMiddleware(), handlers.GetCurrentUser)
<<<<<<< Updated upstream
=======
			auth.GET("/reminders", middleware.AuthMiddleware(), handlers.GetMyReminders)
		}

		// ── Teams ────────────────────────────────────────────────────────────
		teams := api.Group("/teams")
		{
			teams.GET("", handlers.GetTeams)
			teams.GET("/:id", handlers.GetTeamByID)
			teams.POST("/:id/follow", middleware.AuthMiddleware(), handlers.FollowTeam)
			teams.DELETE("/:id/follow", middleware.AuthMiddleware(), handlers.UnfollowTeam)
>>>>>>> Stashed changes
		}
	}

	log.Println("🚀 Nitrous API server starting on :8080")
	r.Run(":8080")
}
