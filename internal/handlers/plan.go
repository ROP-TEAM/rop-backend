package handlers

import (
	dto "ROP_Backend/internal/dto/request"
	"ROP_Backend/internal/dto/response"
	"ROP_Backend/internal/middleware"
	"ROP_Backend/internal/models"
	"ROP_Backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

type PlanHandler struct {
	service *services.PlanningService
}

func NewPlanHandler(s *services.PlanningService) *PlanHandler {
	return &PlanHandler{service: s}
}

// CreatePlan godoc
// @Summary Create a new plan
// @Tags plan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body dto.CreatePlan true "Plan details"
// @Success 201 {object} response.PlanResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/plans [post]
func (h *PlanHandler) Create(c fiber.Ctx) error {
	var req dto.CreatePlan
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	claims := middleware.GetUser(c)
	if claims == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	plan, err := h.service.Create(claims.UserID, req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "plan created",
		"data":    toPlanResponse(plan),
	})
}

// SolvePlan godoc
// @Summary Trigger solver for a plan
// @Tags plan
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Plan ID"
// @Success 200 {object} response.PlanResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/plans/{id}/solve [post]
func (h *PlanHandler) Solve(c fiber.Ctx) error {
	planID := c.Params("id")
	if planID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "plan id required"})
	}

	plan, err := h.service.Solve(c.Context(), planID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "plan solved",
		"data":    toPlanResponse(plan),
	})
}

// GetPlan godoc
// @Summary Get plan by ID with routes and stops
// @Tags plan
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Plan ID"
// @Success 200 {object} response.PlanResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Router /api/plans/{id} [get]
func (h *PlanHandler) GetByID(c fiber.Ctx) error {
	planID := c.Params("id")
	if planID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "plan id required"})
	}

	plan, err := h.service.GetPlan(planID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "plan not found"})
	}

	return c.JSON(fiber.Map{"data": toPlanResponse(plan)})
}

// ListPlans godoc
// @Summary List all plans for the current user's company
// @Tags plan
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} response.PlanResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/plans [get]
func (h *PlanHandler) List(c fiber.Ctx) error {
	claims := middleware.GetUser(c)
	if claims == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	plans, err := h.service.ListPlans(claims.UserID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	result := make([]response.PlanResponse, len(plans))
	for i, p := range plans {
		result[i] = toPlanResponse(&p)
	}

	return c.JSON(fiber.Map{"data": result})
}

func toPlanResponse(p *models.Plan) response.PlanResponse {
	routes := make([]response.RouteResponse, len(p.Routes))
	for i, r := range p.Routes {
		stops := make([]response.StopResponse, len(r.Stops))
		for j, s := range r.Stops {
			stops[j] = response.StopResponse{
				ID:             s.ID,
				SequenceNumber: s.SequenceNumber,
				OrderID:        s.OrderID,
				ArrivalMin:     s.ArrivalMin,
				DepartMin:      s.DepartMin,
			}
		}
		routes[i] = response.RouteResponse{
			ID:            r.ID,
			VehicleID:     r.VehicleID,
			TotalDistance: r.TotalDistance,
			TotalTime:     r.TotalTime,
			Stops:         stops,
		}
	}

	return response.PlanResponse{
		ID:        p.ID,
		Name:      p.Name,
		PlanDate:  p.PlanDate,
		Status:    p.Status,
		DepotLat:  p.DepotLat,
		DepotLon:  p.DepotLon,
		Routes:    routes,
		CreatedAt: p.CreatedAt,
	}
}
