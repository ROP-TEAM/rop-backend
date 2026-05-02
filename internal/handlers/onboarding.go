package handlers

import (
	"ROP_Backend/internal/middleware"
	"ROP_Backend/internal/models"
	"log"

	"github.com/gofiber/fiber/v3"
)

// Onboarding godoc
// @Summary User onboarding
// @Description Save company and user info (step before OTP)
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body models.OnboardingPayload true "Onboarding data"
// @Success 200 {object} handlers.OnboardingResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 401 {object} handlers.ErrorResponse
// @Router /api/onboarding [post]
func (h *UserHandler) Onboarding(c fiber.Ctx) error {
	var body models.OnboardingPayload

	if err := c.Bind().Body(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	claims := middleware.GetUser(c)
	if claims == nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	userID := claims.UserID

	log.Printf("user id pass claims jaa: %s", userID)

	err := h.service.Onboarding(userID, body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "onboarding saved (waiting for OTP)",
	})
}
