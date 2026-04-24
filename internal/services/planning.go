package services

import (
	"context"

	gmaplib "github.com/ROP-TEAM/rop-algorithm/gmap"
	"github.com/ROP-TEAM/rop-algorithm/model"
	"github.com/ROP-TEAM/rop-algorithm/solver"
	"gorm.io/gorm"
)

type PlanningService struct {
	db     *gorm.DB
	matrix gmaplib.DistanceMatrix
	solver solver.Solver
}

func NewPlanningService(db *gorm.DB, matrix gmaplib.DistanceMatrix, s solver.Solver) *PlanningService {
	return &PlanningService{db: db, matrix: matrix, solver: s}
}

// Solve is a placeholder
func (s *PlanningService) Solve(ctx context.Context, p model.Problem) (model.Solution, error) {
	return s.solver.Solve(ctx, p)
}
