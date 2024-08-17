package rest

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lunn06/smart-toy-backend/internal/database/redis"
	"github.com/lunn06/smart-toy-backend/internal/database/sql"
	"github.com/lunn06/smart-toy-backend/internal/models/requests"
)

// @BasePath /api/auth/

// RefreshTokens godoc
// @Summary refresh user's tokens
// @Schemes application/json
// @Description accept json and refresh user tokens
// @Tags authorization
// @Accept json
// @Produce json
// @Param input body requests.RefreshTokensRequest true "session info"
// @Param output body requests.RefreshTokensResponse true "response info"
// @Success 200 "message: RefreshTokens was successful"
// @Failure 401 "error: Invalid to get refresh token from cookie"
// @Failure 500 "error: Invalid to pop token"
// @Failure 500 "error: Invalid to insert token"
// @Failure 500 "error: Invalid to create token"
// @Router /api/auth/refresh [post]
func RefreshTokens(c *gin.Context) {
	// refreshUuid, err := c.Cookie("refreshToken")

	// if err != nil {
	// 	slog.Error(fmt.Sprintf("RefreshToken() error = %v, can't fetch refresh cookies", err))
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"error": "Invalid to get refresh token from cookie",
	// 	})
	// 	return
	// }

	var body requests.RefreshTokensRequest

	if err := c.Bind(&body); err != nil {
		slog.Error(fmt.Sprintf("RefreshToken() error = %v, can't bind RefreshTokensRequest model", err))
		c.JSON(http.StatusBadRequest, requests.RefreshTokensResponse{
			Error: "Invalid request",
		})
		return
	}

	token, err := redis.PopRefreshToken(body.RefreshToken)
	if err != nil {
		slog.Error(fmt.Sprintf("RefreshToken() error = %v, can't delete or select from db", err))
		c.JSON(http.StatusInternalServerError, requests.RefreshTokensResponse{
			Error: "InternalServerError, please try again later",
		})
		return
	}

	if token.FingerPrint != body.SmartToyFingerPrint {
		slog.Error("RefreshToken() request and session fingerprints are not equal", "error", err)
		c.JSON(http.StatusUnauthorized, requests.RefreshTokensResponse{
			Error: "ARE YOU HACKER?",
		})
		return

	}

	user, err := sql.GetUserById(token.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, requests.RefreshTokensResponse{
			Error: "Invalid to get user",
		})
		return

	}

	accessToken, err := newTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, requests.RefreshTokensResponse{
			Error: "Invalid to create token",
		})
		return
	}

	newRefreshUuid, err := redis.InsertRefreshToken(user.Id, body.SmartToyFingerPrint, RefreshLife)
	if err != nil {
		c.JSON(http.StatusBadRequest, requests.RefreshTokensResponse{
			Error: "Invalid to insert token",
		})
		return
	}

	// jwtCookie := http.Cookie{
	// 	Name:     "refreshToken",
	// 	Domain:   config.CFG.HTTPServer.Address,
	// 	Value:    newRefreshUuid,
	// 	MaxAge:   int(RefreshLife.Seconds()),
	// 	Path:     "/api/auth",
	// 	HttpOnly: true,
	// }

	// c.SetSameSite(http.SameSiteStrictMode)

	// c.SetCookie(
	// 	jwtCookie.Name,
	// 	jwtCookie.Value,
	// 	jwtCookie.MaxAge,
	// 	jwtCookie.Path,
	// 	jwtCookie.Domain,
	// 	jwtCookie.Secure,
	// 	jwtCookie.HttpOnly,
	// )

	c.JSON(http.StatusOK, requests.RefreshTokensResponse{
		Message:      "RefreshToken was successful",
		AccessToken:  accessToken,
		RefreshToken: newRefreshUuid,
	})
}
