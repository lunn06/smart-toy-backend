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

var (
	AccessLife  = time.Duration(int64(time.Minute) * int64(config.CFG.AccessLife))
	RefreshLife = time.Duration(int64(time.Minute) * int64(config.CFG.RefreshLife))
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
// @Param output body requests.LoginResponse true "session info"
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
	if err := c.Bind(&body); err != nil {
		slog.Error("Login() can't read request body", "error", err)
		c.JSON(http.StatusBadRequest, requests.LoginResponse{
			Error: "Failed to read body",
		})
		return
	}
	if len(body.Email) > 255 || body.Email == "" {
		slog.Error("Login() invalid email")
		c.JSON(http.StatusUnprocessableEntity, requests.LoginResponse{
			Error: "Email entered incorrectly, because it exceeds the character limit or backwards",
		})
		return
	}
	if len(body.Password) > 72 || body.Password == "" {
		slog.Error("Login() invalid password")
		c.JSON(http.StatusUnprocessableEntity, requests.LoginResponse{
			Error: "Invalid password size",
		})
		return
	}

	user, err := sql.GetUserByEmail(body.Email)
	if err != nil {
		slog.Error("Authentication() can't GetUserByEmail()", "error", err)
		c.JSON(http.StatusForbidden, requests.LoginResponse{
			Error: "Invalid email or password",
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		slog.Error("Login() uncompared password", "error", err)
		c.JSON(http.StatusForbidden, requests.LoginResponse{
			Error: "Invalid email or password",
		})
		return
	}

	accessToken, err := newTokens(user)
	if err != nil {
		slog.Error("newTokens() can't generate new tokens", "error", err)
		c.JSON(http.StatusInternalServerError, requests.LoginResponse{
			Error: "Invalid to create token",
		})
		return
	}

	refreshUuid, err := redis.InsertRefreshToken(user.Id, body.SmartToyFingerPrint, RefreshLife)
	if err != nil {
		slog.Error("Login() can't insert tokens in redis", "error", err)
		c.JSON(http.StatusBadRequest, requests.LoginResponse{
			Error: "Invalid to insert token",
		})
		return
	}

	// jwtCookie := http.Cookie{
	// 	Name:     "refreshToken",
	// 	Domain:   config.CFG.HTTPServer.Address,
	// 	Value:    refreshUuid,
	// 	MaxAge:   int(RefreshLife.Seconds()),
	// 	Path:     "/api/auth",
	// 	HttpOnly: true,
	// 	Secure:   true,
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

	c.JSON(http.StatusOK, requests.LoginResponse{
		Message:      "Authentication was successful",
		AccessToken:  accessToken,
		RefreshToken: refreshUuid,
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
