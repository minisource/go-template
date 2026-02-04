package api

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/minisource/go-common/http/middleware"
	"github.com/minisource/go-common/logging"
	"github.com/minisource/go-common/metrics"
	validation "github.com/minisource/go-common/validations"
	"github.com/minisource/template_go/api/router"
	"github.com/minisource/template_go/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swagger "github.com/swaggo/fiber-swagger"
	_ "github.com/swaggo/files"
	"github.com/swaggo/swag/example/override/docs"
)

var logger = logging.NewLogger(&config.GetConfig().Logger)

func InitServer(cfg *config.Config) {
	// Create Fiber instance
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	RegisterValidators()
	RegisterPrometheus()

	// Middlewares
	app.Use(middleware.DefaultStructuredLogger(&cfg.Logger)) // basic logger
	app.Use(middleware.Prometheus())
	app.Use(middleware.Cors(cfg.Cors.AllowOrigins))
	app.Use(recover.New())
	// app.Use(limiter.New(limiter.Config{
	// 	Max:        100,             // Customize or read from cfg
	// 	Expiration: 1 * time.Minute, // Customize or read from cfg
	// }))

	// Prometheus metrics
	app.Get("/metrics", monitor.New(monitor.Config{Title: "Prometheus Metrics"}))

	// Register routes and Swagger
	RegisterRoutes(app, cfg)
	RegisterSwagger(app, cfg)

	// Start server
	logger := logging.NewLogger(&cfg.Logger)
	logger.Info(logging.General, logging.Startup, "Started", nil)

	err := app.Listen(fmt.Sprintf(":%s", cfg.Server.InternalPort))
	if err != nil {
		logger.Fatal(logging.General, logging.Startup, err.Error(), nil)
	}
}

func RegisterRoutes(app *fiber.App, cfg *config.Config) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Health
	health := v1.Group("/health")
	router.Health(health)

	// Test (add middleware later)
	test := v1.Group("/test")
	router.TestRouter(test)

	// Users
	users := v1.Group("/auth")
	router.User(users, cfg)

	// files := v1.Group("/files")
	// router.File(files, cfg)

	app.Static("/static", "./uploads")

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
}

func RegisterValidators() {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		err := val.RegisterValidation("mobile", validation.IranianMobileNumberValidator, true)
		if err != nil {
			logger.Error(logging.Validation, logging.Startup, err.Error(), nil)
		}
	}
}

func RegisterSwagger(app *fiber.App, cfg *config.Config) {
	docs.SwaggerInfo.Title = "Your Service API"
	docs.SwaggerInfo.Description = "Your Service API - Update this description"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Server.ExternalPort)
	docs.SwaggerInfo.Schemes = []string{"http"}

	// Register swagger route
	app.Get("/swagger/*", swagger.WrapHandler)
}

func RegisterPrometheus() {
	err := prometheus.Register(metrics.DbCall)
	if err != nil {
		logger.Error(logging.Prometheus, logging.Startup, err.Error(), nil)
	}

	err = prometheus.Register(metrics.HttpDuration)
	if err != nil {
		logger.Error(logging.Prometheus, logging.Startup, err.Error(), nil)
	}
}
