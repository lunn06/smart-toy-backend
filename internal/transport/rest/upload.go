package rest

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lunn06/smart-toy-backend/internal/config"
)

// @BasePath /api/

// Upload godoc
// @Summary upload a JSON
// @Schemes application/json
// @Description accepts file and upload it
// @Tags upload
// @Accept json
// @Produce json
// @Success 200 "message: Uploade was successful"
// @Failure 400 "error: Only JSON file accepted"
// @Router /api/upload [post]
func Upload(c *gin.Context) {
	file, _ := c.FormFile("file")

	if !strings.HasSuffix(file.Filename, ".json") {
		slog.Error("Upload() get not json file")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Only JSON file accepted",
		})
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")

	c.SaveUploadedFile(file, config.CFG.UploadDir+file.Filename)

	c.JSON(http.StatusOK, gin.H{
		"message": "Uploade successful",
	})
}
