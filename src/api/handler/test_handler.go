package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/minisource/go-common/http/helper"
)

type header struct {
	UserId  string
	Browser string
}

type personData struct {
	FirstName    string `json:"first_name" binding:"required,alpha,min=4,max=10"`
	LastName     string `json:"last_name" binding:"required,alpha,min=6,max=20"`
	MobileNumber string `json:"mobile_number" binding:"required,mobile,min=11,max=11"`
}
type TestHandler struct {
}

func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

func (h *TestHandler) Test(c *fiber.Ctx) error {
	resp := helper.GenerateBaseResponse("Test", true, 0)
	return c.Status(http.StatusOK).JSON(resp)
}

func (h *TestHandler) Users(c *fiber.Ctx) error {
	resp := helper.GenerateBaseResponse("Users", true, 0)
	return c.Status(http.StatusOK).JSON(resp)
}

// UserById godoc
// @Summary UserById
// @Description UserById
// @Tags Test
// @Accept  json
// @Produce  json
// @Param id path int true "user id"
// @Success 200 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/test/user/{id} [get]
func (h *TestHandler) UserById(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.Status(http.StatusOK).JSON(helper.GenerateBaseResponse(map[string]interface{}{
		"result": "UserById",
		"id":     id,
	}, true, 0))
}

func (h *TestHandler) UserByUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	return c.Status(http.StatusOK).JSON(helper.GenerateBaseResponse(map[string]interface{}{
		"result":   "UserByUsername",
		"username": username,
	}, true, 0))
}

func (h *TestHandler) Accounts(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(helper.GenerateBaseResponse(map[string]interface{}{
		"result": "Accounts",
	}, true, 0))
}

func (h *TestHandler) AddUser(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.Status(http.StatusOK).JSON(helper.GenerateBaseResponse(map[string]interface{}{
		"result": "AddUser",
		"id":     id,
	}, true, 0))
}

func (h *TestHandler) HeaderBinder1(c *fiber.Ctx) error {
	userId := c.Get("UserId") // Fiber's way to get headers
	return c.Status(http.StatusOK).JSON(helper.GenerateBaseResponse(map[string]interface{}{
		"result": "HeaderBinder1",
		"userId": userId,
	}, true, 0))
}

func (h *TestHandler) HeaderBinder2(c *fiber.Ctx) error {
	header := struct {
		UserId string `json:"UserId"`
	}{}

	// Fiber doesn't bind headers into structs directly, so we do it manually
	header.UserId = c.Get("UserId")

	return c.Status(http.StatusOK).JSON(helper.GenerateBaseResponse(map[string]interface{}{
		"result": "HeaderBinder2",
		"header": header,
	}, true, 0))
}

func (h *TestHandler) QueryBinder1(c *fiber.Ctx) error {
	id := c.Query("id")
	name := c.Query("name")
	return c.Status(http.StatusOK).JSON(helper.GenerateBaseResponse(map[string]interface{}{
		"result": "QueryBinder1",
		"id":     id,
		"name":   name,
	}, true, 0))
}

func (h *TestHandler) QueryBinder2(c *fiber.Ctx) error {
	rawQuery := c.Context().QueryArgs()
	// Extract all "id" params manually
	var ids []string
	rawQuery.VisitAll(func(key, val []byte) {
		if string(key) == "id" {
			ids = append(ids, string(val))
		}
	})

	name := c.Query("name")

	return c.Status(fiber.StatusOK).JSON(helper.GenerateBaseResponse(map[string]interface{}{
		"result": "QueryBinder2",
		"ids":    ids,
		"name":   name,
	}, true, 0))
}


// BodyBinder godoc
// @Summary BodyBinder
// @Description BodyBinder
// @Tags Test
// @Accept  json
// @Produce  json
// @Param id path int true "user id"
// @Param name path string true "user name"
// @Success 200 {object} helper.BaseHttpResponse{validationErrors=any{}} "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/test/binder/uri/{id}/{name} [post]
// @Security AuthBearer
func (h *TestHandler) UriBinder(c *fiber.Ctx) error {
	id := c.Params("id")
	name := c.Params("name")

	return c.Status(http.StatusOK).JSON(helper.GenerateBaseResponse(map[string]interface{}{
		"result": "UriBinder",
		"id":     id,
		"name":   name,
	}, true, 0))
}


// BodyBinder godoc
// @Summary BodyBinder
// @Description BodyBinder
// @Tags Test
// @Accept  json
// @Produce  json
// @Param person body personData true "person data"
// @Success 200 {object} helper.BaseHttpResponse{validationErrors=any{}} "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/test/binder/body [post]
// @Security AuthBearer
func (h *TestHandler) BodyBinder(c *fiber.Ctx) error {
	p := personData{}
	if err := c.BodyParser(&p); err != nil {
		resp := helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err)
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	return c.Status(fiber.StatusOK).JSON(helper.GenerateBaseResponse(map[string]interface{}{
		"result": "BodyBinder",
		"person": p,
	}, true, 0))
}

func (h *TestHandler) FormBinder(c *fiber.Ctx) error {
	p := personData{}
	if err := c.BodyParser(&p); err != nil {
		resp := helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err)
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	return c.Status(fiber.StatusOK).JSON(helper.GenerateBaseResponse(map[string]interface{}{
		"result": "FormBinder",
		"person": p,
	}, true, 0))
}

func (h *TestHandler) FileBinder(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		resp := helper.GenerateBaseResponseWithError(nil, false, helper.ValidationError, err)
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	if err := c.SaveFile(file, file.Filename); err != nil {
		resp := helper.GenerateBaseResponseWithError(nil, false, helper.ValidationError, err)
		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	return c.Status(fiber.StatusOK).JSON(helper.GenerateBaseResponse(map[string]interface{}{
		"result": "FileBinder",
		"file":   file.Filename,
	}, true, 0))
}
