package route

import (
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/auth/delivery/http"
	"github.com/gofiber/fiber/v3"
)

type RouteConfig struct {
	App         *fiber.App
	AuthHandler *http.AuthHandler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
}

func (rc *RouteConfig) SetupGuestRoute() {
	api := rc.App.Group("/api/v1")
	api.Post("/users/register", rc.AuthHandler.Register)
	api.Post("/users/login", rc.AuthHandler.Login)
}
