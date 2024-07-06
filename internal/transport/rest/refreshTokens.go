package rest

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lunn06/smart-toy-backend/internal/config"
	"github.com/lunn06/smart-toy-backend/internal/database/redis"
	"github.com/lunn06/smart-toy-backend/internal/database/sql"
	"github.com/lunn06/smart-toy-backend/internal/models"
)

// @BasePath /auth/api/

// RefreshTokens godoc
// @Summary refresh user's tokens
// @Schemes application/json
// @Description accept json and refresh user refresh and access tokens
// @Tags authorization
// @Accept json
// @Produce json
// @Param input body models.LoginRequest true "account info"
// @Success 200 "message: RefreshTokens was successful"
// @Failure 401 "error: Invalid to get refresh token from cookie"
// @Failure 500 "error: Invalid to pop token"
// @Failure 500 "error: Invalid to insert token"
// @Failure 500 "error: Invalid to create token"
// @Router /api/auth/refresh [post]
func RefreshTokens(c *gin.Context) {
	refreshUuid, err := c.Cookie("refreshToken")

	if err != nil {
		slog.Error(fmt.Sprintf("RefreshToken() error = %v, can't fetch refresh cookies", err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid to get refresh token from cookie",
		})
		return
	}

	var body models.RefreshTokensRequest

	if c.Bind(&body) != nil {
		slog.Error(fmt.Sprintf("RefreshToken() error = %v, can't bind RefreshTokensRequest model", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	token, err := redis.PopRefreshToken(refreshUuid)
	if err != nil {
		slog.Error(fmt.Sprintf("RefreshToken() error = %v, can't delete or select from db", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "InternalServerError, please try again later",
		})
		return
	}

	if token.FingerPrint != body.SmartToyFingerPrint {
		slog.Error(fmt.Sprintf("RefreshToken() error = %v, request and session fingerprints are not equal", err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ARE YOU HACKER?",
		})
		return
		
	}

	user, err := sql.GetUserById(token.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid to get user",
		})
		return

	}

	accessToken, err := newTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid to create token",
		})
		return
	}

	newRefreshUuid, err := redis.InsertRefreshToken(user.Id, body.SmartToyFingerPrint, RefreshLife)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid to insert token",
		})
		return
	}

	jwtCookie := http.Cookie{
		Name:     "refreshToken",
		Domain:   config.CFG.HTTPServer.Address,
		Value:    newRefreshUuid,
		MaxAge:   int(RefreshLife.Seconds()),
		Path:     "/api/auth",
		HttpOnly: true,
	}

	c.SetCookie(
		jwtCookie.Name,
		jwtCookie.Value,
		jwtCookie.MaxAge,
		jwtCookie.Path,
		jwtCookie.Domain,
		jwtCookie.Secure,
		jwtCookie.HttpOnly,
	)

	c.JSON(http.StatusOK, gin.H{
		"message":      "RefreshToken was successful",
		"accessToken":  accessToken,
		"refreshToken": newRefreshUuid,
	})
}
