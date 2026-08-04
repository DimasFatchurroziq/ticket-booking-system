package route

import (
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/user/delivery/http"
	"github.com/gofiber/fiber/v3"
)

type RouteConfig struct {
	App         *fiber.App
	UserHandler *http.UserHandler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	// c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/api/users", c.UserHandler.Register)
	// c.App.Post("/api/users/_login", c.UserHandler.Login)
}

// func (c *RouteConfig) SetupAuthRoute() {
// 	// c.App.Use(c.AuthMiddleware)

// 	c.App.Get("/api/hello", c.HelloController.SayHello)

// 	c.App.Delete("/api/users", c.UserHandler.Logout)
// 	c.App.Patch("/api/users/_current", c.UserHandler.Update)
// 	c.App.Get("/api/users/_current", c.UserHandler.Current)
// }
