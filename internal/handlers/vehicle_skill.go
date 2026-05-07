package handlers

import (
	"ROP_Backend/internal/dto"
	"ROP_Backend/internal/middleware"
	"ROP_Backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

type VehicleSkillHandler struct {
	service *services.VehicleSkillService
}

func NewVehicleSkillHandler(s *services.VehicleSkillService) *VehicleSkillHandler {
	return &VehicleSkillHandler{service: s}
}

// VehicleSkills godoc
// @Summary Assign skills to vehicles
// @Description Replace vehicle skills
// @Tags vehicle-skill
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body dto.GroupVehicleSkill true "Vehicle skills"
// @Success 200 {object} handlers.CreateResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/vehicle/skills [post]
func (h *VehicleSkillHandler) Assign(c fiber.Ctx) error {

	var req dto.GroupVehicleSkill

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

	if err := h.service.Assign(req); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "vehicle skills assigned",
	})
}
