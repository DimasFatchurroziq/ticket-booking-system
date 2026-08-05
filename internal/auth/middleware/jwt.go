package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/DimasFatchurroziq/ticket-booking-system/internal/auth/token"
)

func Auth(tokenManager token.TokenManager) fiber.Handler {
	return func(c fiber.Ctx) error {

		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "authorization header is required",
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid authorization header",
			})
		}

		claims, err := tokenManager.Parse(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid token",
			})
		}

		c.Locals("user", claims)

		return c.Next()
	}
}
