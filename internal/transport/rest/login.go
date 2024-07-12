package rest

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lunn06/smart-toy-backend/internal/config"
	"github.com/lunn06/smart-toy-backend/internal/database/redis"
	"github.com/lunn06/smart-toy-backend/internal/database/sql"
	"github.com/lunn06/smart-toy-backend/internal/models"
	"github.com/lunn06/smart-toy-backend/internal/models/requests"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessLife  = time.Minute * 30
	RefreshLife = time.Hour * 24 * 30
)

// @BasePath /api/auth/

// Login godoc
// @Summary login the user
// @Schemes application/json
// @Description accepts json with user info and authorize him
// @Tags authorization
// @Accept json
// @Produce json
// @Param input body requests.LoginRequest true "account info"
// @Success 200 "message: Login was successful"
// @Failure 400 "error: Failed to read body"
// @Failure 422 "error: Email entered incorrectly, because it exceeds the character limit or backwards"
// @Failure 422 "error: Invalid password size"
// @Failure 403 "error: Invalid email or password"
// @Failure 500 "error: Invalid to create token"
// @Failure 400 "error: Invalid to insert token"
// @Router /api/auth/login [post]
func Login(c *gin.Context) {
	body := requests.LoginRequest{}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	if len(body.Email) > 255 || body.Email == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Email entered incorrectly, because it exceeds the character limit or backwards",
		})
		return
	}
	if len(body.Password) > 72 || body.Password == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Invalid password size",
		})
		return
	}
	user, err := sql.GetUserByEmail(body.Email)
	if err != nil {
		slog.Error("Authentication() can't GetUserByEmail()", "error", err)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Invalid email or password",
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

	refreshUuid, err := redis.InsertRefreshToken(user.Id, body.SmartToyFingerPrint, RefreshLife)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid to insert token",
		})
		return
	}

	jwtCookie := http.Cookie{
		Name:     "refreshToken",
		Domain:   config.CFG.HTTPServer.Address,
		Value:    refreshUuid,
		MaxAge:   int(RefreshLife.Seconds()),
		Path:     "/api/auth",
		HttpOnly: true,
		Secure: true,
	}

	c.SetSameSite(http.SameSiteStrictMode)

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
		"message":      "Authentication was successful",
		"accessToken":  accessToken,
		"refreshToken": refreshUuid,
	})
}

func newTokens(user models.User) (string, error) {
	var jwtSecretKey = []byte(config.CFG.JWTSecretKey)

	accessPayload := jwt.MapClaims{
		"sub":   user.Id,
		"email": user.Email,
		"exp":   AccessLife,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessPayload)
	signedAccessToken, err := accessToken.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return signedAccessToken, nil
}
