package http

import (
	"errors"

	"github.com/DimasFatchurroziq/ticket-booking-system/internal/auth/token"
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/user/domain"
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/user/usecase"
	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
	validate    *validator.Validate
	Log         *zap.Logger
}

func NewUserHandler(app *fiber.App, userUsecase usecase.UserUsecase, Log *zap.Logger, validate *validator.Validate) *UserHandler {

	return &UserHandler{
		userUsecase: userUsecase,
		validate:    validate,
		Log:         Log,
	}
}

func (h *UserHandler) Me(c fiber.Ctx) error {

	claims, ok := c.Locals("user").(*token.Claims)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	cmd := usecase.MeCommand{
		UserID: claims.UserID,
	}
	result, err := h.userUsecase.Me(c.Context(), cmd)

	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "credential tidak valid",
			})
		}

		return c.Status(500).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.JSON(result)
}
