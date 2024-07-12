package rest

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lunn06/smart-toy-backend/internal/database/sql"
	"github.com/lunn06/smart-toy-backend/internal/models"
	"github.com/lunn06/smart-toy-backend/internal/models/requests"
	"golang.org/x/crypto/bcrypt"
)

// @BasePath /api/auth/

// Registration godoc
// @Summary register user
// @Schemes application/json
// @Description accepts json with user info and registers him
// @Tags authorization
// @Accept json
// @Produce json
// @Param input body requests.RegisterRequest true "account info"
// @Success 200 "message: Registration was successful"
// @Failure 400 "error: Failed to read body"
// @Failure 422 "error: Failed create email, because it exceeds the character limit or backwards"
// @Failure 422 "error: Failed create channel_name, because it exceeds the character limit or backwards"
// @Failure 422 "error: Failed create password, because it exceeds the character limit or backwards"
// @Failure 500 "error: Failed to hash password. Please, try again later"
// @Failure 409 "error: email or channel already been use"
// @Router /api/auth/registration [post]
func Registration(c *gin.Context) {
	body := requests.RegisterRequest{}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	if len(body.Email) > 255 || body.Email == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Failed create email, because it exceeds the character limit or backwards",
		})
		return
	}
	if len(body.Password) > 72 || body.Password == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Failed create password, because it exceeds the character limit or backwards",
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error on the server. Please, try again later",
		})
		return
	}
	user := models.User{
		Email:    body.Email,
		Password: string(hash),
	}
	err = sql.InsertUser(user)
	slog.Error("Registration() can't InsertUser", "error", err)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "email or channel already been use",
		})
		return
	}

	// accessToken, refreshToken, err := newTokens(user)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "Invalid to create token",
	// 	})
	// 	return
	// }

	// refreshUUID, err := sql.InsertToken(userID, refreshToken)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "Invalid to insert token",
	// 	})
	// 	return
	// }

	// jwtCookie := http.Cookie{
	// 	Name:     "refreshToken",
	// 	Value:    refreshUUID,
	// 	MaxAge:   refreshLife,
	// 	Path:     "/api/auth",
	// 	HttpOnly: true,
	// }

	// c.SetCookie(
	// 	jwtCookie.Name,
	// 	jwtCookie.Value,
	// 	jwtCookie.MaxAge,
	// 	jwtCookie.Path,
	// 	jwtCookie.Domain,
	// 	jwtCookie.Secure,
	// 	jwtCookie.HttpOnly,
	// )

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration was successful",
		// "accessToken":  accessToken,
		// "refreshToken": refreshUUID,
	})
}
