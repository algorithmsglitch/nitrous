package database

import "nitrous-backend/models"

// Add these to your existing var block in db.go:
//
// var (
//     ...existing vars...
//     Teams     []models.Team
//     Streams   []models.Stream
//     Reminders []models.Reminder
//     Orders    []models.Order
//     Follows   []models.TeamFollow
// )
//
// And inside InitDB(), add:
//     Teams   = SeedTeams()
//     Streams = SeedStreams()

// These are declared here to avoid touching db.go directly.
// If you prefer, move them into db.go alongside Events, Users, etc.
var (
	Teams     []models.Team
	Streams   []models.Stream
	Reminders []models.Reminder
	Orders    []models.Order
	Follows   []models.TeamFollow
)

// InitNewCollections seeds the new data. Call this from InitDB() in db.go.
func InitNewCollections() {
	Teams   = SeedTeams()
	Streams = SeedStreams()
}
