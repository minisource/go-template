package config

import (
	"errors"
	"log"
	"os"

	auth "github.com/minisource/auth/service"
	gormdb "github.com/minisource/go-common/db/gorm"
	"github.com/minisource/go-common/http/middleware"
	"github.com/minisource/go-common/logging"
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	Gorm   gormdb.GormConfig
	Cors   CorsConfig
	Logger logging.LoggerConfig
	Auth   auth.AuthServiceConfig
	OTP    middleware.OtpConfig
}

type ServerConfig struct {
	InternalPort            string
	ExternalPort            string
	RunMode                 string
	Domain                  string
	RefreshCookieMaxAgeSecs int // Max age for refresh token cookie in seconds (default: 604800 = 7 days)
}

type CorsConfig struct {
	AllowOrigins string
}

func GetConfig() *Config {
	cfgPath := getConfigPath(os.Getenv("APP_ENV"))
	v, err := LoadConfig(cfgPath, "yml")
	if err != nil {
		log.Fatalf("Error in load config %v", err)
	}

	cfg, err := ParseConfig(v)
	envPort := os.Getenv("PORT")
	if envPort != "" {
		cfg.Server.ExternalPort = envPort
		log.Printf("Set external port from environment -> %s", cfg.Server.ExternalPort)
	} else {
		cfg.Server.ExternalPort = cfg.Server.InternalPort
		log.Printf("Set external port from environment -> %s", cfg.Server.ExternalPort)
	}
	if err != nil {
		log.Fatalf("Error in parse config %v", err)
	}

	return cfg
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var cfg Config
	err := v.Unmarshal(&cfg)
	if err != nil {
		log.Printf("Unable to parse config: %v", err)
		return nil, err
	}
	return &cfg, nil
}
func LoadConfig(filename string, fileType string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType(fileType)
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		log.Printf("Unable to read config: %v", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}
	return v, nil
}

func getConfigPath(env string) string {
	if env == "docker" {
		return "/app/config/config-docker"
	} else if env == "production" {
		return "/config/config-production"
	} else {
		return "../config/config-development"
	}
}
