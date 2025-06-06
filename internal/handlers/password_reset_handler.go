package handlers

import (
	"e_meeting/internal/models"
	"e_meeting/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type PasswordResetHandler struct {
	passwordResetService *services.PasswordResetService
	validator            *validator.Validate
}

func NewPasswordResetHandler(passwordResetService *services.PasswordResetService) *PasswordResetHandler {
	return &PasswordResetHandler{
		passwordResetService: passwordResetService,
		validator:            validator.New(),
	}
}

func (h *PasswordResetHandler) RequestReset(c *fiber.Ctx) error {
	var req models.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Call service to handle reset request

	link, err := h.passwordResetService.RequestReset(c.Context(), req.Email, c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to process password reset request")
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process password reset request",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "If your email is registered, you will receive a password reset link",
		"reset_link": link,
	})
}

func (h *PasswordResetHandler) ResetPassword(c *fiber.Ctx) error {

	var req models.ResetPasswordConfirmRequest
	err := c.QueryParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	log.Info().Msgf("Reset password request: %v", req)
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Call service to handle password reset
	if err := h.passwordResetService.ResetPassword(c.Context(), req.Token, req.NewPassword); err != nil {
		log.Error().Err(err).Msg("Failed to reset password")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to reset password",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password has been reset successfully",
	})
}
