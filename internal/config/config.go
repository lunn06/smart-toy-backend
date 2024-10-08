package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

var CFG Config

type Config struct {
	HTTPServer `yaml:"http_server"`
	Database   `yaml:"database"`
	Redis      `yaml:"redis"`
	JWT        `yaml:"jwt"`
	SmartToy   `yaml:"smart_toy"`
	UploadDir  string `yaml:"upload_dir"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:"127.0.0.1"`
	Port    string `yaml:"port" env-default:"8080"`
}

type Database struct {
	Driver   string `yaml:"driver"`
	Address  string `yaml:"address"`
	User     string `yaml:"user"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"ssl_mode"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
}

type JWT struct {
	AccessLife   int    `yaml:"access_life"`
	RefreshLife  int    `yaml:"refresh_life"`
	JWTSecretKey string `yaml:"jwt_secret_key"`
}

type SmartToy struct {
	MaxCount int `yaml:"max_count"`
}

func Init() {
	CFG = MustLoad("configs/main.yaml")
}

func MustLoad(configPath string) Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return cfg
}
