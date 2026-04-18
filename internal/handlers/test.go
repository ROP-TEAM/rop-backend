package handlers

import (
	"ROP_Backend/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

// GetMe godoc
// @Summary Get current user
// @Description Get user_id from JWT
// @Tags user
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} handlers.UserResponse
// @Failure 401 {object} handlers.ErrorResponse
// @Router /api/test [get]
func Test(c fiber.Ctx) error {
	user := middleware.GetUser(c)

	if user == nil {
		return c.Status(401).JSON(ErrorResponse{
			Error: "Invalid token",
		})
	}

	return c.JSON(UserResponse{
		UserID: user.UserID,
		Email:  user.Email,
	})
}

// TestOTP godoc
// @Summary Test OTP endpoint
// @Description Use for testing OTP rate limiter (1 request per minute)
// @Tags otp
// @Success 200 {object} handlers.OTPResponse
// @Failure 429 {object} handlers.ErrorResponse
// @Router /api/auth/otp [post]
func TestOTP(c fiber.Ctx) error {
	return c.JSON(OTPResponse{
		Message: "OTP endpoint hit",
	})
}
