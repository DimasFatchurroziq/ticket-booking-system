package config

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
)

func NewFiber(config *viper.Viper) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      config.GetString("APP_NAME"),
		ErrorHandler: NewErrorHandler(),
	})

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx fiber.Ctx, err error) error { // 👈 Di Fiber v3, ctx adalah fiber.Ctx (bukan *fiber.Ctx)
		// Default status code 500
		code := fiber.StatusInternalServerError

		// Pengecekan aman menggunakan errors.As (Best Practice di Go)
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}

		// Response format JSON
		return ctx.Status(code).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}
}
