package handlers

import (
	"ROP_Backend/internal/services"

	"gorm.io/gorm"
)

type PlanHandler struct {
	planService *services.PlanService
}

func NewPlanHandler(db *gorm.DB) *PlanHandler {
	planService := services.NewPlanService(db)
	return &PlanHandler{
		planService: planService,
	}
}

// func Create(c fiber.Ctx) error {
// 	var body dto.CreatePlanRequest
// 	if err := c.Bind().Body(&body); err != nil {
// 		// return c.JSON()
// 	}
// }
