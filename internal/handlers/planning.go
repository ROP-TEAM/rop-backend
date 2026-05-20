package handlers

import (
	"context"
	"time"

	dto "ROP_Backend/internal/dto/request"
	"ROP_Backend/internal/dto/response"

	"github.com/gofiber/fiber/v3"
)

type planningOptimizer interface {
	Optimize(ctx context.Context, req dto.OptimizeRequest) (*response.OptimizeResponse, error)
}

type PlanningHandler struct {
	planningService planningOptimizer
}

func NewPlanningHandler(s planningOptimizer) *PlanningHandler {
	return &PlanningHandler{planningService: s}
}

func statusToHTTP(status string) int {
	switch status {
	case "INFEASIBLE":
		return 422
	case "TIMEOUT":
		return 504
	default:
		return 200
	}
}

// Optimize godoc
// @Summary Optimize vehicle routing
// @Description Optimize vehicle routing problem given vehicles and orders
// @Tags planning
// @Accept json
// @Produce json
// @Param body body dto.OptimizeRequest true "optimize request"
// @Success 200 {object} response.OptimizeResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 422 {object} response.OptimizeResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Failure 504 {object} handlers.ErrorResponse
// @Router /api/optimize [post]
func (h *PlanningHandler) Optimize(c fiber.Ctx) error {
	var req dto.OptimizeRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	result, err := h.planningService.Optimize(ctx, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(statusToHTTP(result.Status)).JSON(fiber.Map{
		"message": "optimize success",
		"data":    result,
	})
}
