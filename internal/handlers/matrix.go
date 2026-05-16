package handlers

import (
	dto "ROP_Backend/internal/dto/request"
	"ROP_Backend/internal/dto/response"
	"ROP_Backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

type MatrixHandler struct {
	service *services.MatrixService
}

func NewMatrixHandler(s *services.MatrixService) *MatrixHandler {
	return &MatrixHandler{service: s}
}

// BuildMatrix godoc
// @Summary Build distance matrix
// @Description Build an n×n distance matrix from a list of lat/lng locations. Durations in minutes, distances in meters.
// @Tags matrix
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body dto.MatrixRequest true "List of locations (min 2)"
// @Success 200 {object} response.BuildMatrixResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/matrix [post]
func (h *MatrixHandler) BuildMatrix(c fiber.Ctx) error {
	var req dto.MatrixRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{Error: "invalid body"})
	}

	if len(req.Locations) < 2 {
		return c.Status(400).JSON(ErrorResponse{Error: "at least 2 locations required"})
	}

	result, err := h.service.BuildMatrix(c.Context(), req.Locations)
	if err != nil {
		return c.Status(500).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.JSON(response.BuildMatrixResponse{
		Message: "matrix built",
		Node:    len(req.Locations),
		Result:  result,
	})
}
