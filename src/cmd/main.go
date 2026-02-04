package main

import (
	"github.com/minisource/template_go/api"
	"github.com/minisource/template_go/config"
	"github.com/minisource/template_go/infra/persistence/migration"
	auth "github.com/minisource/auth/service"
	"github.com/minisource/go-common/db/gorm"
	"github.com/minisource/go-common/logging"
)

// @securityDefinitions.apikey AuthBearer
// @in header
// @name Authorization
func main() {
	cfg := config.GetConfig()
	logger := logging.NewLogger(&cfg.Logger)

	auth := auth.NewAuthService(cfg.Auth)
	err := auth.HealthCheck()
	if err != nil {
		logger.Fatal(logging.Casdoor, logging.Startup, err.Error(), nil)
	}

	err = gormdb.InitDb(&cfg.Gorm)
	defer gormdb.CloseDb()
	if err != nil {
		logger.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}
	migration.Up1()

	api.InitServer(cfg)
}
