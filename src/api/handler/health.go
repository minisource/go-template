package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/minisource/go-common/http/helper"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck godoc
// @Summary Health Check
// @Description Health Check
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/health/ [get]
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	resp := helper.GenerateBaseResponse("Working!", true, 0)
	return c.Status(http.StatusOK).JSON(resp)
}
