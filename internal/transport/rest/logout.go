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
// @Param input body requests.LogoutRequest true "session info"
// @Param output body requests.LogoutResponse true "response info"
// @Success 200 "message: Logout was successful"
// @Failure 400 "error: Failed to read body"
// @Failure 500 "error: Invalid to remove session"
// @Router /api/auth/logout [delete]
func Logout(c *gin.Context) {
	body := requests.LogoutRequest{}
	if err := c.Bind(&body); err != nil {
		slog.Error("Logout() can't PopRefreshToken", "error", err)
		c.JSON(http.StatusBadRequest, requests.LogoutResponse{
			Error: "Failed to read body",
		})
		return
	}

	_, err := redis.PopRefreshToken(body.RefreshToken)

	if err != nil {
		slog.Error("Logout() can't PopRefreshToken", "error", err)
		c.JSON(http.StatusInternalServerError, requests.LogoutResponse{
			Error: "Invalid to remove session",
		})
		return
	}

	c.JSON(http.StatusOK, requests.LogoutResponse{
		Message: "Logout was successful",
	})
}
