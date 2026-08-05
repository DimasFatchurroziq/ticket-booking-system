package route

import (
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/auth/middleware"
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/auth/token"
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/user/delivery/http"
	"github.com/gofiber/fiber/v3"
)

type RouteConfig struct {
	App          *fiber.App
	UserHandler  *http.UserHandler
	TokenManager token.TokenManager
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
}

func (rc *RouteConfig) SetupGuestRoute() {
	api := rc.App.Group("/api/v1")
	api.Use(middleware.Auth(rc.TokenManager))
	api.Get("/users/me", rc.UserHandler.Me)
}
