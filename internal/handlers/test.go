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
