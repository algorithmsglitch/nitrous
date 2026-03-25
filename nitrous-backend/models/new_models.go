package models

import "time"

// ── TEAM ──────────────────────────────────────────────────────────────────────

type Team struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`  // e.g. "MOTORSPORT · F1"
	Country   string    `json:"country"`   // e.g. "Austria"
	Founded   int       `json:"founded"`
	Rank      int       `json:"rank"`
	Wins      int       `json:"wins"`
	Points    int       `json:"points"`
	Following int       `json:"following"`
	Drivers   []string  `json:"drivers"`
	Color     string    `json:"color"`     // accent color key
	CreatedAt time.Time `json:"createdAt"`
}

type TeamFollow struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	TeamID    string    `json:"teamId"`
	CreatedAt time.Time `json:"createdAt"`
}

// ── STREAM ────────────────────────────────────────────────────────────────────

type Stream struct {
	ID            string    `json:"id"`
	EventID       string    `json:"eventId"`
	Title         string    `json:"title"`
	Subtitle      string    `json:"subtitle"`       // e.g. "Lap 87 / 200"
	Category      string    `json:"category"`
	Location      string    `json:"location"`
	Quality       string    `json:"quality"`        // "4K" | "HD"
	Viewers       int       `json:"viewers"`
	IsLive        bool      `json:"isLive"`
	CurrentLeader string    `json:"currentLeader"`
	CurrentSpeed  string    `json:"currentSpeed"`
	Color         string    `json:"color"`
	CreatedAt     time.Time `json:"createdAt"`
}

// StreamTelemetry is broadcast over WebSocket to update live data
type StreamTelemetry struct {
	StreamID      string `json:"streamId"`
	Viewers       int    `json:"viewers"`
	CurrentLeader string `json:"currentLeader"`
	CurrentSpeed  string `json:"currentSpeed"`
	Subtitle      string `json:"subtitle"`
}

// ── REMINDER ──────────────────────────────────────────────────────────────────

type Reminder struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	EventID   string    `json:"eventId"`
	CreatedAt time.Time `json:"createdAt"`
}

// ── ORDER (MERCH) ─────────────────────────────────────────────────────────────

type OrderItem struct {
	MerchID  string `json:"merchId"`
	Name     string `json:"name"`
	Price    float64 `json:"price"`
	Quantity int    `json:"quantity"`
	Size     string `json:"size,omitempty"`
}

type Order struct {
	ID        string      `json:"id"`
	UserID    string      `json:"userId"`
	Items     []OrderItem `json:"items"`
	Total     float64     `json:"total"`
	Status    string      `json:"status"` // "pending" | "confirmed" | "shipped"
	CreatedAt time.Time   `json:"createdAt"`
}

type CreateOrderRequest struct {
	Items []OrderItem `json:"items" binding:"required,min=1"`
}
