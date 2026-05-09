package handlers

import (
	dto "ROP_Backend/internal/dto/request"
	"ROP_Backend/internal/middleware"
	"ROP_Backend/internal/services"

	"strconv"

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
// @Success 200 {object} response.VehicleGroupResponse
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

	data, err := h.service.GroupCreate(req)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "vehicles created",
		"data":    data,
	})
}

// PatchVehicle godoc
// @Summary Update vehicle
// @Description Patch vehicle by id
// @Tags vehicle
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Vehicle ID"
// @Param body body dto.UpdateVehicle true "Vehicle update"
// @Success 200 {object} response.VehicleGroupResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/vehicles/{id} [patch]
func (h *VehicleHandler) Patch(c fiber.Ctx) error {

	var req dto.UpdateVehicle

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid vehicle id",
		})
	}

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

	data, err := h.service.Update(
		uint(id),
		req,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "vehicle updated",
		"data":    data,
	})
}

// DeleteVehicle godoc
// @Summary Delete vehicle list
// @Description Delete multiple vehicles
// @Tags vehicle
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body dto.DeleteVehicle true "Vehicle IDs"
// @Success 200 {object} handlers.CreateResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/vehicles [delete]
func (h *VehicleHandler) Delete(c fiber.Ctx) error {

	var req dto.DeleteVehicle

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	if len(req.ID) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "id required",
		})
	}

	claims := middleware.GetUser(c)
	if claims == nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	if err := h.service.Delete(req); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "vehicles deleted",
	})
}
