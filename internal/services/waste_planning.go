package services

import (
	"context"

	dto "ROP_Backend/internal/dto/request"
	"ROP_Backend/internal/dto/response"

	"github.com/ROP-TEAM/rop-algorithm/solver"
	"github.com/ROP-TEAM/rop-algorithm/waste"
)

// WastePlanningService plans a week of garbage collection by mapping cleaned
// route units onto the shared VRPTW solver. Phase 1 builds its own haversine
// matrix inside the waste package, so it needs no distance-matrix provider; an
// OSRM matrix is a later refinement.
type WastePlanningService struct {
	solver solver.Solver
}

func NewWastePlanningService(solver solver.Solver) *WastePlanningService {
	return &WastePlanningService{solver: solver}
}

// Plan runs the full week and returns the metrics next to the current-assignment
// baseline.
func (s *WastePlanningService) Plan(
	ctx context.Context,
	req *dto.WastePlanRequest,
) (*response.WastePlanResponse, error) {
	planner := s.buildPlanner(req)

	report, err := planner.PlanWeek(ctx, req.Units)
	if err != nil {
		return nil, err
	}
	return toWasteResponse(report), nil
}

func (s *WastePlanningService) buildPlanner(req *dto.WastePlanRequest) waste.Planner {
	cfg := waste.DefaultConfig()
	if req.CollectKmh > 0 {
		cfg.CollectKmh = req.CollectKmh
	}
	if req.DriveKmh > 0 {
		cfg.DriveKmh = req.DriveKmh
	}
	cfg.TimeLimitMS = req.TimeLimitMS

	mode := waste.MinimizeFleet
	if req.Mode == "full" {
		mode = waste.FullFleet
	}

	return waste.Planner{
		Solver:      s.solver,
		Config:      cfg,
		Fleet:       waste.Fleet{Mode: mode, Count: req.Count, ShiftStart: req.ShiftStart, ShiftEnd: req.ShiftEnd},
		MaxPerCell:  req.MaxPerCell,
		MaxShiftMin: req.MaxShiftMin,
	}
}

func toWasteResponse(report waste.WeekReport) *response.WastePlanResponse {
	perDay := make([]response.WasteDayMetric, 0, len(report.PerDay))
	for _, day := range report.PerDay {
		perDay = append(perDay, response.WasteDayMetric{
			Weekday:             day.Weekday,
			Cells:               day.Cells,
			Bands:               day.Bands,
			TrucksUsed:          day.Metrics.TrucksUsed,
			AssignedUnits:       day.Metrics.AssignedUnits,
			UnassignedUnits:     day.Metrics.UnassignedUnits,
			InterStopDeadheadKm: km(day.Metrics.InterStopDeadheadM),
			MeanUtilisation:     day.Metrics.MeanUtilisation,
			LoadGini:            day.Metrics.LoadGini,
		})
	}

	t := report.Totals
	return &response.WastePlanResponse{
		Baseline: response.WasteBaseline{
			TruckGroups:         report.Baseline.TruckGroups,
			AssignedUnits:       report.Baseline.AssignedUnits,
			InterStopDeadheadKm: km(report.Baseline.InterStopDeadheadM),
		},
		PerDay: perDay,
		Totals: response.WasteTotals{
			TruckShifts:           t.TruckShifts,
			AssignedUnits:         t.AssignedUnits,
			UnassignedUnits:       t.UnassignedUnits,
			InterStopDeadheadKm:   km(t.InterStopDeadheadM),
			InterStopReductionPct: t.InterStopReductionPct,
		},
	}
}

func km(metres float64) float64 { return float64(int(metres/100+0.5)) / 10 }
