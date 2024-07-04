package app

import (
	"log"

	"github.com/lunn06/smart-toy-backend/internal/config"
	"github.com/lunn06/smart-toy-backend/internal/services"
	"github.com/lunn06/video-hosting/internal/database"
)

func init() {
	config.Init()
	database.Init()
	//initializers.ParseConfig()
	//initializers.ConnectToDB()
}

func Run() {
	r := services.SetupRouter()

	onlyPortAddress := ":" + config.CFG.HTTPServer.Port
	log.Fatal(r.Run(onlyPortAddress))
}
