package handlers

import (
	"context"

	"ROP_Backend/internal/config"
	"ROP_Backend/internal/services"

	"github.com/ROP-TEAM/rop-algorithm/model"
	"github.com/gofiber/fiber/v3"
)

type MatrixHandler struct {
	cfg *config.Config
}

func NewMatrixHandler(cfg *config.Config) *MatrixHandler {
	return &MatrixHandler{cfg: cfg}
}

type matrixTestRequest struct {
	Locations []locationInput `json:"locations"`
}

type locationInput struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type matrixTestResponse struct {
	Durations [][]int `json:"durations"`
	Distances [][]int `json:"distances"`
}

// TestMatrix godoc
// @Summary Test distance matrix
// @Description Build an n×n distance matrix from a list of lat/lng locations. Durations in minutes, distances in meters.
// @Tags matrix
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body matrixTestRequest true "List of locations (min 2)"
// @Success 200 {object} matrixTestResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/matrix/test [post]
func (h *MatrixHandler) Test(c fiber.Ctx) error {
	var req matrixTestRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{Error: "invalid body"})
	}

	if len(req.Locations) < 2 {
		return c.Status(400).JSON(ErrorResponse{Error: "at least 2 locations required"})
	}

	dm, err := services.NewDistanceMatrix(h.cfg)
	if err != nil {
		return c.Status(500).JSON(ErrorResponse{Error: err.Error()})
	}

	locs := make([]model.Location, len(req.Locations))
	for i, l := range req.Locations {
		locs[i] = model.NewLatLngLocation(l.Lat, l.Lng)
	}

	durations, distances, err := dm.BuildMatrix(context.Background(), locs, model.MatrixOptions{})
	if err != nil {
		return c.Status(500).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.JSON(matrixTestResponse{
		Durations: durations,
		Distances: distances,
	})
}