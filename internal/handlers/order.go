package handlers

import (
	dto "ROP_Backend/internal/dto/request"
	"ROP_Backend/internal/middleware"
	"ROP_Backend/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(s *services.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}

// CreateOrder godoc
// @Summary Create order list
// @Description Create multiple orders
// @Tags order
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body dto.GroupCreateOrder true "Order list"
// @Success 200 {object} response.OrderResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/orders [post]
func (h *OrderHandler) Create(c fiber.Ctx) error {
	var req dto.GroupCreateOrder

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
		"message": "orders created",
		"data":    data,
	})
}

// PatchOrder godoc
// @Summary Update order
// @Description Patch order by id
// @Tags order
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Order ID"
// @Param body body dto.UpdateOrder true "Order update"
// @Success 200 {object} response.OrderResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/orders/{id} [patch]
func (h *OrderHandler) Patch(c fiber.Ctx) error {

	var req dto.UpdateOrder

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid order id",
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
		"message": "order updated",
		"data":    data,
	})
}

// DeleteOrder godoc
// @Summary Delete order list
// @Description Delete multiple orders
// @Tags order
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body dto.DeleteVehicle true "Order IDs"
// @Success 200 {object} handlers.CreateResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/orders [delete]
func (h *OrderHandler) Delete(c fiber.Ctx) error {

	var req dto.DeleteOrder

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

	// claims := middleware.GetUser(c)
	// if claims == nil {
	// 	return c.Status(401).JSON(fiber.Map{
	// 		"error": "unauthorized",
	// 	})
	// }

	if err := h.service.Delete(req); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "orders deleted",
	})
}
