package handlers

import (
	dto "ROP_Backend/internal/dto/request"
	"ROP_Backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

type MatrixHandler struct {
	service *services.MatrixService
}

func NewMatrixHandler(s *services.MatrixService) *MatrixHandler {
	return &MatrixHandler{service: s}
}

// TestMatrix godoc
// @Summary Test distance matrix
// @Description Build an n×n distance matrix from a list of lat/lng locations. Durations in minutes, distances in meters.
// @Tags matrix
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body dto.MatrixTestRequest true "List of locations (min 2)"
// @Success 200 {object} matrixTestResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/matrix/test [post]
func (h *MatrixHandler) Test(c fiber.Ctx) error {
	var req dto.MatrixTestRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{Error: "invalid body"})
	}

	if len(req.Locations) < 2 {
		return c.Status(400).JSON(ErrorResponse{Error: "at least 2 locations required"})
	}

	result, err := h.service.BuildMatrix(req.Locations)
	if err != nil {
		return c.Status(500).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "matrix Built ",
		"node":    len(req.Locations),
		"result":  result,
	})
}
