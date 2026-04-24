package handlers

import (
	"ROP_Backend/internal/dto"

	"github.com/gofiber/fiber/v3"
)

// CreatePlan godoc
// @Summary Create plan
// @Description Receive cars and orders from frontend (no processing yet)
// @Tags plan
// @Accept json
// @Produce json
// @Param body body dto.PlanRequest true "Plan payload"
// @Success 200 {object} handlers.PlanResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Router /api/v1/plan [post]
func CreatePlan(c fiber.Ctx) error {
	var body dto.PlanRequest

	if err := c.Bind().Body(&body); err != nil {
		return c.Status(400).JSON(ErrorResponse{
			Error: "invalid body",
		})
	}
	return c.Status(fiber.StatusOK).JSON(PlanResponse{
		Message: "Plan created successfully",
	})
}
