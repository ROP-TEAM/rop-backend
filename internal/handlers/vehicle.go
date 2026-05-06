package handlers

import (
	"ROP_Backend/internal/dto"
	"ROP_Backend/internal/middleware"
	"ROP_Backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

type VehicleHandler struct {
	service *services.VehicleService
}

func NewVehicleHandler(s *services.VehicleService) *VehicleHandler {
	return &VehicleHandler{service: s}
}

// CreateVehicle godoc
// @Summary Create vehicle list
// @Description Create multiple vehicles
// @Tags vehicle
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body dto.GroupCreateVehicle true "Vehicle list"
// @Success 200 {object} handlers.CreateResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/vehicles [post]
func (h *VehicleHandler) Create(c fiber.Ctx) error {
	var req dto.GroupCreateVehicle

	if err := c.Bind().Body(&req); err != nil {
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

	if err := h.service.GroupCreate(req); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "vehicels created",
	})

}
