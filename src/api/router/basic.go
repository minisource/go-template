package router

import (
	"github.com/minisource/template_go/api/handler"
	"github.com/minisource/template_go/config"
	"github.com/gofiber/fiber/v2"
)

const GetByFilterExp string = "/get-by-filter"

func File(r fiber.Router, cfg *config.Config) {
	h := handler.NewFileHandler(cfg)

	r.Post("/", h.Create)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
	r.Get("/:id", h.GetById)
	r.Post(GetByFilterExp, h.GetByFilter)
}
