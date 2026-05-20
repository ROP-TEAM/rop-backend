package handlers

import (
	"context"
	"time"

	request "ROP_Backend/internal/dto/request"
	"ROP_Backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

type PlanningHandler struct {
	planningService *services.PlanningService
}

func NewPlanningHandler(
	planningService *services.PlanningService,
) *PlanningHandler {
	return &PlanningHandler{
		planningService: planningService,
	}
}

// Optimize godoc
// @Summary Optimize vehicle routing
// @Description Optimize vehicle routing problem given vehicles and orders
// @Tags planning
// @Accept json
// @Produce json
// @Param body body request.OptimizeRequest true "optimize request"
// @Success 200 {object} response.OptimizeResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 401 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/optimize [post]
func (h *PlanningHandler) Optimize(
	c fiber.Ctx,
) error {

	var req request.OptimizeRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(
		c.Context(),
		10*time.Second,
	)
	defer cancel()

	result, err :=
		h.planningService.Optimize(
			ctx,
			req,
		)

	if err != nil {
		return c.Status(
			fiber.StatusInternalServerError,
		).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(
		fiber.StatusOK,
	).JSON(fiber.Map{
		"message": "optimize success",
		"data":    result,
	})
}
