package handlers

import (
	dto "ROP_Backend/internal/dto/request"
	"ROP_Backend/internal/middleware"
	"ROP_Backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

type TagSkillHandler struct {
	service *services.TagSkillService
}

func NewTagSkillHandler(s *services.TagSkillService) *TagSkillHandler {
	return &TagSkillHandler{service: s}
}

// CreateTagSkill godoc
// @Summary Create skill list
// @Description Create multiple skills by plan
// @Tags skill
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body dto.GroupCreateTagSkill true "Skill list"
// @Success 200 {object} response.TagSkillGroupResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/skills [post]
func (h *TagSkillHandler) Create(c fiber.Ctx) error {
	var req dto.GroupCreateTagSkill

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
		"message": "skills updated",
		"data":    data,
	})
}
