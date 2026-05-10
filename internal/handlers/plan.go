package handlers

import (
	"ROP_Backend/internal/dto"
	"ROP_Backend/internal/middleware"
	"ROP_Backend/internal/services"
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"
)

type PlanHandler struct {
	service *services.PlanService
}

func NewPlanHandler(service *services.PlanService) *PlanHandler {
	return &PlanHandler{service: service}
}

// CreatePlan godoc
// @Summary Create plan
// @Description Create a plan
// @Tags plan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body dto.CreatePlanRequest true "plan"
// @Success 200 {object} dto.CreatePlanResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/plan [post]
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

// PatchPlan godoc
// @Summary Patch plan name by plan id
// @Description Patch a plan name by id by attach plan id via url params, and updated name by body
// @Tags plan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body dto.UpdatePlanNameByIDRequest true "plan"
// @Success 200 {object} dto.UpdatePlanNameByIDResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/plan/:id [patch]
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

// DeletePlan godoc
// @Summary Delete plan by plan id
// @Description hard delete plan by id by attach plan id via url params, it will delte all its legacy
// @Tags plan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body dto.DeletePlanByIDRequest true "plan"
// @Success 200 {object} dto.DeletePlanByIDResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/plan/:id [delete]
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

// DuplicatePlanByID godoc
// @Summary DuplicatePlan and its legacy by id
// @Description user has to be in the same company as target plan to allow duplicating it
// @Tags plan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body dto.DuplicatePlanByIDRequest true "plan"
// @Success 200 {object} dto.DuplicatePlanByIDResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/plan/:id [post]
func (h *PlanHandler) DuplicateByID(c fiber.Ctx) error {
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

	res, err := h.service.DuplicateByID(ctx, id, claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound),
			errors.Is(err, services.ErrPlanNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})

		case errors.Is(err, services.ErrCreatingOrder),
			errors.Is(err, services.ErrCreatingVehicle),
			errors.Is(err, services.ErrCreatingTagSkill),
			errors.Is(err, services.ErrCreatingOrderSkill),
			errors.Is(err, services.ErrCreatingVehicleSkill):
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "internal server error",
				"details": err.Error(),
			})

		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "plan duplicated",
		"data":    res,
	})
}
