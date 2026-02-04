package router

import (
	"github.com/minisource/template_go/api/handler"
	"github.com/minisource/template_go/config"
	"github.com/gofiber/fiber/v2"
)

func User(r fiber.Router, cfg *config.Config) {
	h := handler.NewUserHandler(cfg)

	r.Post("/send-otp" /*, middleware.OtpLimiter(&cfg.OTP)*/, h.SendOtp)
	r.Post("/login-by-mobile", h.RegisterLoginByMobileNumber)
	// r.Post("/refresh-token", h.RefreshToken)
}
