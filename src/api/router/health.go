package router

import (
	"github.com/minisource/template_go/api/handler"
	"github.com/gofiber/fiber/v2"
)

func Health(r fiber.Router) {
	h := handler.NewHealthHandler()
	r.Get("/", h.Health)
}
