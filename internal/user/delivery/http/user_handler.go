package http

import (
	"errors"

	"github.com/DimasFatchurroziq/ticket-booking-system/internal/user/domain"
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/user/dto"
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

func (h *UserHandler) Register(c fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "validation failed",
		})
	}

	cmd := usecase.RegisterCommand{
		Email:       req.Email,
		Password:    req.Password,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
	}

	result, err := h.userUsecase.Register(c.Context(), cmd)
	if err != nil {

		if errors.Is(err, domain.ErrEmailExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "email already exists",
			})
		}

		return c.Status(500).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.RegisterResponse{
		ID:      result.ID,
		Email:   result.Email,
		Message: "registrasi berhasil",
	})
}
