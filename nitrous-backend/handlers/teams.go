package handlers

import (
	"net/http"
	"nitrous-backend/database"
<<<<<<< Updated upstream

	"github.com/gin-gonic/gin"
=======
	"nitrous-backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
>>>>>>> Stashed changes
)

// GetTeams returns all teams
func GetTeams(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"teams": database.Teams,
		"count": len(database.Teams),
	})
}

<<<<<<< Updated upstream
// GetTeamByID returns a single team by ID
=======
// GetTeamByID returns a single team
>>>>>>> Stashed changes
func GetTeamByID(c *gin.Context) {
	id := c.Param("id")

	for _, team := range database.Teams {
		if team.ID == id {
			c.JSON(http.StatusOK, team)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
}

<<<<<<< Updated upstream
// FollowTeam adds the authenticated user to the team's followers
func FollowTeam(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")

	for i, team := range database.Teams {
		if team.ID == id {
			// check if already following
			uid := userID.(string)
			for _, f := range team.Followers {
				if f == uid {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Already following"})
					return
				}
			}

			database.Teams[i].Followers = append(database.Teams[i].Followers, uid)
			database.Teams[i].FollowersCount = len(database.Teams[i].Followers)

			c.JSON(http.StatusOK, gin.H{"message": "Team followed", "team": database.Teams[i]})
=======
// FollowTeam adds a follow record for the authenticated user
func FollowTeam(c *gin.Context) {
	teamID := c.Param("id")
	userID := c.GetString("userID")

	// Check team exists
	teamIdx := -1
	for i, t := range database.Teams {
		if t.ID == teamID {
			teamIdx = i
			break
		}
	}
	if teamIdx == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	// Check not already following
	for _, f := range database.Follows {
		if f.UserID == userID && f.TeamID == teamID {
			c.JSON(http.StatusConflict, gin.H{"error": "Already following this team"})
>>>>>>> Stashed changes
			return
		}
	}

<<<<<<< Updated upstream
	c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
}

// UnfollowTeam removes the authenticated user from the team's followers
func UnfollowTeam(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")
	uid := userID.(string)

	for i, team := range database.Teams {
		if team.ID == id {
			// find follower index
			idx := -1
			for j, f := range team.Followers {
				if f == uid {
					idx = j
					break
				}
			}

			if idx == -1 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Not following"})
				return
			}

			// remove follower
			database.Teams[i].Followers = append(database.Teams[i].Followers[:idx], database.Teams[i].Followers[idx+1:]...)
			database.Teams[i].FollowersCount = len(database.Teams[i].Followers)

			c.JSON(http.StatusOK, gin.H{"message": "Team unfollowed", "team": database.Teams[i]})
=======
	follow := models.TeamFollow{
		ID:        uuid.New().String(),
		UserID:    userID,
		TeamID:    teamID,
		CreatedAt: time.Now(),
	}

	database.Follows = append(database.Follows, follow)
	database.Teams[teamIdx].Following++

	c.JSON(http.StatusOK, gin.H{
		"message":   "Now following team",
		"following": database.Teams[teamIdx].Following,
	})
}

// UnfollowTeam removes a follow record
func UnfollowTeam(c *gin.Context) {
	teamID := c.Param("id")
	userID := c.GetString("userID")

	teamIdx := -1
	for i, t := range database.Teams {
		if t.ID == teamID {
			teamIdx = i
			break
		}
	}
	if teamIdx == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	for i, f := range database.Follows {
		if f.UserID == userID && f.TeamID == teamID {
			database.Follows = append(database.Follows[:i], database.Follows[i+1:]...)
			if database.Teams[teamIdx].Following > 0 {
				database.Teams[teamIdx].Following--
			}
			c.JSON(http.StatusOK, gin.H{
				"message":   "Unfollowed team",
				"following": database.Teams[teamIdx].Following,
			})
>>>>>>> Stashed changes
			return
		}
	}

<<<<<<< Updated upstream
	c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
=======
	c.JSON(http.StatusNotFound, gin.H{"error": "Follow record not found"})
>>>>>>> Stashed changes
}
