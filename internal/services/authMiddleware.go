package services

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lunn06/smart-toy-backend/internal/database/redis"
	"github.com/lunn06/smart-toy-backend/internal/database/sql"
	"github.com/lunn06/smart-toy-backend/internal/transport/rest"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshUUID, err := c.Cookie("refreshToken")

		if err != nil {
			slog.Error("AuthMiddleware() can't get refreshToken cookie", "error", err)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("AuthMiddleware() error = %v", err),
			})
			return
		}

		token, err := redis.GetRefreshToken(refreshUUID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "INVALID_REFRESH_SESSION: refresh token out of life",
			})
			return
		}

		if token.CreationTime.Add(time.Duration(rest.RefreshLife)).Compare(time.Now()) == 1 {
			slog.Error("AuthMiddleware() error = refresh token out of life")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "INVALID_REFRESH_SESSION: refresh token out of life",
			})
			return
		}

		if _, err = sql.GetUserById(token.UserId); err != nil {
			slog.Error("AuthMiddleware() can't GetUserById", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "INVALID_REFRESH_SESSION: no user with this token",
			})
			return
		}

		c.Next()
	}
}
