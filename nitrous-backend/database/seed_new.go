package database

// Add these variables alongside your existing ones in db.go:
//
//   var (
//       Teams     []models.Team
//       Streams   []models.Stream
//       Reminders []models.Reminder
//       Orders    []models.Order
//       Follows   []models.TeamFollow
//   )
//
// And call seedTeams() and seedStreams() inside InitDB().
// This file contains those seed functions.

import (
	"nitrous-backend/models"
	"time"

	"github.com/google/uuid"
)

func SeedTeams() []models.Team {
	return []models.Team{
		{
			ID:        uuid.New().String(),
			Name:      "Red Bull Racing",
			Category:  "MOTORSPORT · F1",
			Country:   "Austria",
			Founded:   2005,
			Rank:      1,
			Wins:      21,
			Points:    860,
			Following: 8200000,
			Drivers:   []string{"Max Verstappen", "Sergio Pérez"},
			Color:     "red",
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			Name:      "Hendrick Motorsports",
			Category:  "MOTORSPORT · NASCAR",
			Country:   "USA",
			Founded:   1984,
			Rank:      2,
			Wins:      14,
			Points:    2340,
			Following: 3100000,
			Drivers:   []string{"Kyle Larson", "Chase Elliott", "William Byron", "Alex Bowman"},
			Color:     "cyan",
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			Name:      "Toyota Gazoo Racing",
			Category:  "RALLY · WRC",
			Country:   "Japan",
			Founded:   1957,
			Rank:      3,
			Wins:      9,
			Points:    564,
			Following: 1800000,
			Drivers:   []string{"Sébastien Ogier", "Elfyn Evans", "Kalle Rovanperä"},
			Color:     "orange",
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			Name:      "Team Sea Force",
			Category:  "WATER · SPEED BOAT",
			Country:   "Italy",
			Founded:   2010,
			Rank:      4,
			Wins:      7,
			Points:    320,
			Following: 420000,
			Drivers:   []string{"F. Bertrand", "L. Capelli"},
			Color:     "blue",
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			Name:      "Falcon Air Squadron",
			Category:  "AIR · AIR RACING",
			Country:   "France",
			Founded:   2015,
			Rank:      5,
			Wins:      5,
			Points:    198,
			Following: 280000,
			Drivers:   []string{"A. Garnier", "B. Morin"},
			Color:     "purple",
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			Name:      "Baja Iron Squad",
			Category:  "OFF-ROAD · TROPHY TRUCK",
			Country:   "Mexico",
			Founded:   1998,
			Rank:      6,
			Wins:      12,
			Points:    415,
			Following: 640000,
			Drivers:   []string{"C. Wedekin", "P. McMillin"},
			Color:     "gold",
			CreatedAt: time.Now(),
		},
	}
}

func SeedStreams() []models.Stream {
	return []models.Stream{
		{
			ID:            uuid.New().String(),
			EventID:       "", // link to actual event ID at runtime if needed
			Title:         "NASCAR Daytona 500",
			Subtitle:      "Lap 87 / 200",
			Category:      "MOTORSPORT",
			Location:      "Daytona International Speedway · FL",
			Quality:       "4K",
			Viewers:       1200000,
			IsLive:        true,
			CurrentLeader: "Bubba Wallace #23",
			CurrentSpeed:  "198 mph",
			Color:         "red",
			CreatedAt:     time.Now(),
		},
		{
			ID:            uuid.New().String(),
			EventID:       "",
			Title:         "World Dirt Track Championship",
			Subtitle:      "Heat 3 — Semi Finals",
			Category:      "MOTORSPORT",
			Location:      "Knob Noster · Missouri, USA",
			Quality:       "HD",
			Viewers:       340000,
			IsLive:        true,
			CurrentLeader: "Kyle Larson #57",
			CurrentSpeed:  "142 mph",
			Color:         "orange",
			CreatedAt:     time.Now(),
		},
		{
			ID:            uuid.New().String(),
			EventID:       "",
			Title:         "Lake Como Speed Boat Qualifier",
			Subtitle:      "Qualifying Round 2",
			Category:      "WATER",
			Location:      "Lake Como · Italy",
			Quality:       "HD",
			Viewers:       89000,
			IsLive:        true,
			CurrentLeader: "F. Bertrand #9",
			CurrentSpeed:  "87 knots",
			Color:         "cyan",
			CreatedAt:     time.Now(),
		},
		{
			ID:            uuid.New().String(),
			EventID:       "",
			Title:         "Red Bull Skydive Series",
			Subtitle:      "Live Drop — 14,800ft",
			Category:      "AIR",
			Location:      "Interlaken Drop Zone · Switzerland",
			Quality:       "HD",
			Viewers:       220000,
			IsLive:        true,
			CurrentLeader: "A. Garnier",
			CurrentSpeed:  "120 mph",
			Color:         "purple",
			CreatedAt:     time.Now(),
		},
	}
}
