package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minisource/go-common/http/helper"
	"github.com/minisource/template_go/api/dto"
	"github.com/minisource/template_go/config"
	"github.com/minisource/template_go/constant"
	"github.com/minisource/template_go/dependency"
	"github.com/minisource/template_go/usecase"
)

type UsersHandler struct {
	userUsecase *usecase.UserUsecase
	config      *config.Config
}

func NewUserHandler(cfg *config.Config) *UsersHandler {
	userUsecase := usecase.NewUserUsecase(cfg, dependency.GetUserRepository(cfg))
	return &UsersHandler{userUsecase: userUsecase, config: cfg}
}

// SendOtp godoc
// @Summary Send otp to user
// @Description Send otp to user
// @Tags Users
// @Accept  json
// @Produce  json
// @Param Request body dto.GetOtpRequest true "GetOtpRequest"
// @Success 201 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Failure 409 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/users/send-otp [post]
func (h *UsersHandler) SendOtp(c *fiber.Ctx) error {
	req := new(dto.GetOtpRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err),
		)
	}

	if err := h.userUsecase.SendOtpByMobileNumber(req.CountryCode, req.MobileNumber); err != nil {
		return c.Status(helper.TranslateErrorToStatusCode(err)).JSON(
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		helper.GenerateBaseResponse(nil, true, helper.Success),
	)
}

// RegisterLoginByMobileNumber godoc
// @Summary RegisterLoginByMobileNumber
// @Description RegisterLoginByMobileNumber
// @Tags Users
// @Accept  json
// @Produce  json
// @Param Request body dto.RegisterLoginByMobileRequest true "RegisterLoginByMobileRequest"
// @Success 201 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Failure 409 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/users/login-by-mobile [post]
func (h *UsersHandler) RegisterLoginByMobileNumber(c *fiber.Ctx) error {
	req := new(dto.RegisterLoginByMobileRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err),
		)
	}

	token, err := h.userUsecase.RegisterAndLoginByMobileNumber(c.Context(), req.CountryCode, req.MobileNumber, req.Otp)
	if err != nil {
		return c.Status(helper.TranslateErrorToStatusCode(err)).JSON(
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err),
		)
	}

	// Set the refresh token in a cookie
	maxAge := h.config.Server.RefreshCookieMaxAgeSecs
	if maxAge == 0 {
		maxAge = 604800 // Default: 7 days in seconds
	}
	c.Cookie(&fiber.Cookie{
		Name:     constant.RefreshTokenCookieName,
		Value:    token.RefreshToken,
		MaxAge:   maxAge,
		Path:     "/",
		Domain:   h.config.Server.Domain,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Strict",
	})

	return c.Status(fiber.StatusOK).JSON(
		helper.GenerateBaseResponse(token, true, helper.Success),
	)
}
