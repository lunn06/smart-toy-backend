package app

import (
	"log"

	"github.com/lunn06/smart-toy-backend/internal/config"
	"github.com/lunn06/smart-toy-backend/internal/database/redis"
	"github.com/lunn06/smart-toy-backend/internal/database/sql"
	"github.com/lunn06/smart-toy-backend/internal/services"
)

func init() {
	config.Init()
	sql.Init()
	redis.Init()
}

func Run() {
	r := services.SetupRouter()

	onlyPortAddress := ":" + config.CFG.HTTPServer.Port
	log.Fatal(r.Run(onlyPortAddress))
}
