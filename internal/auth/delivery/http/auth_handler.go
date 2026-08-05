package http

import (
	"errors"

	"github.com/DimasFatchurroziq/ticket-booking-system/internal/auth/domain"
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/auth/dto"
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/auth/usecase"
	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	validate    *validator.Validate
	Log         *zap.Logger
}

func NewAuthHandler(app *fiber.App, authUsecase usecase.AuthUsecase, Log *zap.Logger, validate *validator.Validate) *AuthHandler {

	return &AuthHandler{
		authUsecase: authUsecase,
		validate:    validate,
		Log:         Log,
	}
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
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

	result, err := h.authUsecase.Register(c.Context(), cmd)
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

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req dto.LoginRequest

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

	cmd := usecase.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.authUsecase.Login(c.Context(), cmd)
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

	return c.Status(fiber.StatusCreated).JSON(dto.LoginResponse{
		Token:   result.AccessToken,
		Message: "login berhasil",
	})
}
