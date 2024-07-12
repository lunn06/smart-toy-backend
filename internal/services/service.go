package services

import (
	"github.com/gin-gonic/gin"
	"github.com/lunn06/smart-toy-backend/internal/transport/rest"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// r.MaxMultipartMemory = 8 << 20  // 8 MiB

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/registration", rest.Registration)
			auth.POST("/login", rest.Login)
			auth.POST("/refresh", rest.RefreshTokens)

			auth.DELETE("/logout", rest.Logout)
		}

		api.POST("/upload", rest.Upload)
	}

	SetupDocs(r)

	return r
}