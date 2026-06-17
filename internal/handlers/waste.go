package handlers

import (
	"context"
	"time"

	dto "ROP_Backend/internal/dto/request"
	"ROP_Backend/internal/dto/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type wastePlanner interface {
	Plan(ctx context.Context, req *dto.WastePlanRequest) (*response.WastePlanResponse, error)
}

type WasteHandler struct {
	wasteService wastePlanner
}

func NewWasteHandler(s wastePlanner) *WasteHandler {
	return &WasteHandler{wasteService: s}
}

// Plan godoc
// @Summary Plan a week of garbage collection
// @Description Plan collection routes for a week from cleaned route units and compare deadhead to the current assignment
// @Tags waste
// @Accept json
// @Produce json
// @Param body body dto.WastePlanRequest true "waste plan request"
// @Success 200 {object} response.WastePlanResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/waste/plan [post]
func (h *WasteHandler) Plan(c fiber.Ctx) error {
	var req dto.WastePlanRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := validator.New().Struct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "validation failed",
			"details": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Minute)
	defer cancel()

	result, err := h.wasteService.Plan(ctx, &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "waste plan success",
		"data":    result,
	})
}
