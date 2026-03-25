package handlers

import (
	"encoding/json"
	"log"
	"net/http"
<<<<<<< Updated upstream
=======
	"nitrous-backend/database"
	"nitrous-backend/models"
	"strconv"
>>>>>>> Stashed changes
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

<<<<<<< Updated upstream
type Stream struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	IsLive      bool   `json:"isLive"`
	ViewerCount int    `json:"viewerCount"`
}

type TelemetryPayload struct {
	Type      string    `json:"type"`
	StreamID  string    `json:"streamId"`
	SpeedKPH  int       `json:"speedKph"`
	RPM       int       `json:"rpm"`
	Gear      int       `json:"gear"`
	GForce    float64   `json:"gForce"`
	Timestamp time.Time `json:"timestamp"`
}

type Hub struct {
	clients    map[*websocket.Conn]bool
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	broadcast  chan []byte
}

var (
	streams = []Stream{
		{ID: "stream-1", Title: "Daytona 500 — Main Feed", Category: "motorsport", IsLive: true, ViewerCount: 12042},
		{ID: "stream-2", Title: "Dakar Rally — Stage Cam", Category: "offroad", IsLive: true, ViewerCount: 5421},
		{ID: "stream-3", Title: "Sky Racing Cockpit View", Category: "air", IsLive: false, ViewerCount: 0},
	}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	hubOnce   sync.Once
	streamHub = &Hub{
		clients:    make(map[*websocket.Conn]bool),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		broadcast:  make(chan []byte, 128),
	}
)

func ensureHubRunning() {
	hubOnce.Do(func() {
		go RunHub()
	})
}

// GetStreams returns all stream feeds.
func GetStreams(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"streams": streams,
		"count":   len(streams),
	})
}

// GetStreamByID returns one stream feed.
func GetStreamByID(c *gin.Context) {
	id := c.Param("id")

	for _, stream := range streams {
		if stream.ID == id {
			c.JSON(http.StatusOK, stream)
=======
// ── REST ──────────────────────────────────────────────────────────────────────

// GetStreams returns all live streams
func GetStreams(c *gin.Context) {
	var liveStreams []models.Stream
	for _, s := range database.Streams {
		if s.IsLive {
			liveStreams = append(liveStreams, s)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"streams": liveStreams,
		"count":   len(liveStreams),
	})
}

// GetStreamByID returns a single stream
func GetStreamByID(c *gin.Context) {
	id := c.Param("id")

	for _, s := range database.Streams {
		if s.ID == id {
			c.JSON(http.StatusOK, s)
>>>>>>> Stashed changes
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Stream not found"})
}

<<<<<<< Updated upstream
// StreamsWS upgrades the request to websocket and registers the client to telemetry updates.
func StreamsWS(c *gin.Context) {
	ensureHubRunning()

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to establish websocket connection"})
		return
	}

	streamHub.register <- conn

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			streamHub.unregister <- conn
=======
// ── WEBSOCKET ─────────────────────────────────────────────────────────────────

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow localhost:3000 and production frontend
		//origin := r.Header.Get("Origin")
		//return origin == "http://localhost:3000" || origin == "https://nitrous.vercel.app"
		return true
	},
}

// hub holds all active WebSocket connections
type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	mu        sync.Mutex
}

var streamHub = &Hub{
	clients:   make(map[*websocket.Conn]bool),
	broadcast: make(chan []byte, 256),
}

// RunHub starts the broadcast loop — call this once as a goroutine in main.go:
//
//	go handlers.RunHub()
func RunHub() {
	for msg := range streamHub.broadcast {
		streamHub.mu.Lock()
		for conn := range streamHub.clients {
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				conn.Close()
				delete(streamHub.clients, conn)
			}
		}
		streamHub.mu.Unlock()
	}
}

// StreamsWS upgrades the connection and registers the client
// Route: GET /ws/streams
func StreamsWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	streamHub.mu.Lock()
	streamHub.clients[conn] = true
	streamHub.mu.Unlock()

	// Send current state immediately on connect
	snapshot, _ := json.Marshal(database.Streams)
	conn.WriteMessage(websocket.TextMessage, snapshot)

	// Keep connection alive; remove on disconnect
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			streamHub.mu.Lock()
			delete(streamHub.clients, conn)
			streamHub.mu.Unlock()
>>>>>>> Stashed changes
			break
		}
	}
}

<<<<<<< Updated upstream
// RunHub runs the websocket client registration, unregistration, and broadcast loop.
func RunHub() {
	for {
		select {
		case conn := <-streamHub.register:
			streamHub.clients[conn] = true

		case conn := <-streamHub.unregister:
			if _, ok := streamHub.clients[conn]; ok {
				delete(streamHub.clients, conn)
				_ = conn.Close()
			}

		case message := <-streamHub.broadcast:
			for conn := range streamHub.clients {
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					delete(streamHub.clients, conn)
					_ = conn.Close()
				}
			}
		}
	}
}

// BroadcastTelemetry publishes telemetry updates to all connected websocket clients.
func BroadcastTelemetry(streamID string, speedKPH int, rpm int, gear int, gForce float64) {
	ensureHubRunning()

	payload := TelemetryPayload{
		Type:      "telemetry",
		StreamID:  streamID,
		SpeedKPH:  speedKPH,
		RPM:       rpm,
		Gear:      gear,
		GForce:    gForce,
		Timestamp: time.Now().UTC(),
	}

	message, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal telemetry payload: %v", err)
		return
	}

	select {
	case streamHub.broadcast <- message:
	default:
		log.Printf("telemetry message dropped: broadcast channel is full")
	}
}
=======
// BroadcastTelemetry pushes a telemetry update to all connected clients.
// Call this from a ticker or an admin endpoint when race data changes.
func BroadcastTelemetry(t models.StreamTelemetry) {
	// Update in-memory store
	for i, s := range database.Streams {
		if s.ID == t.StreamID {
			database.Streams[i].Viewers = t.Viewers
			database.Streams[i].CurrentLeader = t.CurrentLeader
			database.Streams[i].CurrentSpeed = t.CurrentSpeed
			database.Streams[i].Subtitle = t.Subtitle
			break
		}
	}

	payload, err := json.Marshal(t)
	if err != nil {
		return
	}
	streamHub.broadcast <- payload
}

// SimulateTelemetry is a dev helper that sends fake updates every 5s.
// Call as a goroutine in main.go during development:
//
//	go handlers.SimulateTelemetry()
func SimulateTelemetry() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	lap := 87
	for range ticker.C {
		if len(database.Streams) == 0 {
			continue
		}
		lap++
		s := database.Streams[0]
		BroadcastTelemetry(models.StreamTelemetry{
			StreamID:      s.ID,
			Viewers:       s.Viewers + 100,
			CurrentLeader: s.CurrentLeader,
			CurrentSpeed:  s.CurrentSpeed,
			Subtitle:      "Lap " + itoa(lap) + " / 200",
		})
	}
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
>>>>>>> Stashed changes
