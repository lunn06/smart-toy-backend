package rest

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lunn06/smart-toy-backend/internal/database/redis"
	"github.com/lunn06/smart-toy-backend/internal/models/requests"
)

// @BasePath /api/auth/

// Logout godoc
// @Summary delete user session
// @Schemes application/json
// @Description accepts json with refresh token and delete session
// @Tags authorization
// @Accept json
// @Produce json
// @Param input body requests.LogoutRequest true "refresh token"
// @Success 200 "message: Logout was successful"
// @Failure 400 "error: Failed to read body"
// @Failure 500 "error: Invalid to remove session"
// @Router /api/auth/logout [delete]
func Logout(c *gin.Context) {
	body := requests.LogoutRequest{}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	err := redis.DelRefreshToken(body.RefreshToken)

	if err != nil {
		slog.Error("Logout() can't PopRefreshToken", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid to remove session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout was successful",
	})
}
