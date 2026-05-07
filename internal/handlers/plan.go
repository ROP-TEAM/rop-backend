package handlers

import (
	"ROP_Backend/internal/dto"
	"ROP_Backend/internal/middleware"
	"ROP_Backend/internal/services"
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
)

type PlanHandler struct {
	service *services.PlanService
}

func NewPlanHandler(service *services.PlanService) *PlanHandler {
	return &PlanHandler{service: service}
}

func (h *PlanHandler) Create(c fiber.Ctx) error {
	var body dto.CreatePlanRequest
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	if body.Name == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "name field is required",
		})
	}

	body.SetDefaults()

	claims := middleware.GetUser(c)
	if claims == nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	userID := claims.UserID

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	res, err := h.service.CreateByUserID(ctx, userID, &body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "plan created",
		"data":    res,
	})
}

func (h *PlanHandler) UpdateNameByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	var body dto.UpdatePlanNameByIDRequest
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	body.SetDefaults()

	claims := middleware.GetUser(c)
	if claims == nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}
	//ยังไม่เซฟเท่าไหร่ เพราะว่าขอแค่มีtoken แล้วรู้ไอดีของplanก็ใช้ได้เลยย

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	res, err := h.service.UpdateNameByID(ctx, id, claims.UserID, &body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "name updated",
		"data":    res,
	})
}

func (h *PlanHandler) DeleteByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	claims := middleware.GetUser(c)
	if claims == nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	err := h.service.DeleteByID(ctx, id, claims.UserID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "plan deleted",
		// "data":    res,
	})
}
