package router

import (
	"github.com/minisource/template_go/api/handler"
	"github.com/gofiber/fiber/v2"
)

func TestRouter(r fiber.Router) {
	h := handler.NewTestHandler()

	r.Get("/", h.Test)
	r.Get("/users", h.Users)
	r.Get("/user/:id", h.UserById)
	r.Get("/user/get-user-by-username/:username", h.UserByUsername)
	r.Get("/user/:id/accounts", h.Accounts)
	r.Post("/add-user", h.AddUser)

	r.Post("/binder/header1", h.HeaderBinder1)
	r.Post("/binder/header2", h.HeaderBinder2)

	r.Post("/binder/query1", h.QueryBinder1)
	r.Post("/binder/query2", h.QueryBinder2)

	r.Post("/binder/uri/:id/:name", h.UriBinder)
	r.Post("/binder/body", h.BodyBinder)
	r.Post("/binder/form", h.FormBinder)
	r.Post("/binder/file", h.FileBinder)
}
