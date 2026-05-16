package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	dto "ROP_Backend/internal/dto/request"
	"ROP_Backend/internal/models"
	"ROP_Backend/internal/repository"

	"github.com/ROP-TEAM/rop-algorithm/model"
	"github.com/ROP-TEAM/rop-algorithm/solver"
	"gorm.io/gorm"
)

type PlanningService struct {
	planRepo    *repository.PlanRepository
	orderRepo   *repository.OrderRepository
	vehicleRepo *repository.VehicleRepository
	userRepo    *repository.UserRepository
	matrixSvc   *MatrixService
	solver      solver.Solver
}

func NewPlanningService(
	db *gorm.DB,
	matrixSvc *MatrixService,
	slvr solver.Solver,
) *PlanningService {
	return &PlanningService{
		planRepo:    repository.NewPlanRepository(db),
		orderRepo:   repository.NewOrderRepository(db),
		vehicleRepo: repository.NewVehicleRepository(db),
		userRepo:    repository.NewUserRepository(db),
		matrixSvc:   matrixSvc,
		solver:      slvr,
	}
}

func (s *PlanningService) Create(userID uint, req dto.CreatePlan) (*models.Plan, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if user.CompanyID == nil {
		return nil, errors.New("user has no company — complete onboarding first")
	}

	planDate, err := time.Parse("2006-01-02", req.PlanDate)
	if err != nil {
		return nil, fmt.Errorf("invalid plan_date format, expected YYYY-MM-DD: %w", err)
	}

	plan := &models.Plan{
		Name:      req.Name,
		PlanDate:  planDate,
		DepotLat:  req.DepotLat,
		DepotLon:  req.DepotLon,
		Status:    "pending",
		CompanyID: *user.CompanyID,
	}

	if err := s.planRepo.Create(plan); err != nil {
		return nil, err
	}
	return plan, nil
}

func (s *PlanningService) Solve(ctx context.Context, planID string) (*models.Plan, error) {
	plan, err := s.planRepo.FindByID(planID)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}

	orders, err := s.orderRepo.FindByPlanWithSkills(planID)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, errors.New("plan has no orders")
	}

	vehicles, err := s.vehicleRepo.FindByPlanWithSkills(planID)
	if err != nil {
		return nil, err
	}
	if len(vehicles) == 0 {
		return nil, errors.New("plan has no vehicles")
	}

	if err := s.planRepo.UpdateStatus(planID, "processing"); err != nil {
		return nil, err
	}

	locations := buildLocations(plan, orders)
	matrixResp, err := s.matrixSvc.BuildMatrix(locations)
	if err != nil {
		s.planRepo.UpdateStatus(planID, "failed")
		return nil, fmt.Errorf("matrix build failed: %w", err)
	}

	problem := buildProblem(plan, orders, vehicles, intsToFloats(matrixResp.Durations), intsToFloats(matrixResp.Distances))

	solution, err := s.solver.Solve(ctx, problem)
	if err != nil {
		s.planRepo.UpdateStatus(planID, "failed")
		return nil, fmt.Errorf("solver failed: %w", err)
	}

	routes := buildRoutes(planID, orders, vehicles, solution)
	if len(routes) > 0 {
		if err := s.planRepo.SaveRoutes(routes); err != nil {
			s.planRepo.UpdateStatus(planID, "failed")
			return nil, err
		}
	}

	status := "completed"
	if solution.Status == model.SolutionStatusInfeasible {
		status = "failed"
	}
	if err := s.planRepo.UpdateStatus(planID, status); err != nil {
		return nil, err
	}

	return s.planRepo.FindByIDWithRoutes(planID)
}

func (s *PlanningService) GetPlan(planID string) (*models.Plan, error) {
	return s.planRepo.FindByIDWithRoutes(planID)
}

func (s *PlanningService) ListPlans(userID uint) ([]models.Plan, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if user.CompanyID == nil {
		return nil, errors.New("user has no company")
	}
	return s.planRepo.FindByCompanyID(*user.CompanyID)
}

// --- helpers ---

func buildLocations(plan *models.Plan, orders []models.Order) []dto.LocationInput {
	locs := make([]dto.LocationInput, 0, 1+len(orders))
	locs = append(locs, dto.LocationInput{Lat: plan.DepotLat, Lng: plan.DepotLon})
	for _, o := range orders {
		lat, lng := 0.0, 0.0
		if o.DesLatitude != nil {
			lat = *o.DesLatitude
		}
		if o.DesLongitude != nil {
			lng = *o.DesLongitude
		}
		locs = append(locs, dto.LocationInput{Lat: lat, Lng: lng})
	}
	return locs
}

func buildProblem(
	plan *models.Plan,
	orders []models.Order,
	vehicles []models.Vehicle,
	durations [][]float64,
	distances [][]float64,
) model.Problem {
	depot := model.Node{
		ID:   "depot",
		Lat:  plan.DepotLat,
		Lng:  plan.DepotLon,
		Type: model.NodeTypeDepot,
	}

	nodes := make([]model.Node, len(orders))
	for i, o := range orders {
		nodes[i] = orderToNode(o)
	}

	modelVehicles := make([]model.Vehicle, len(vehicles))
	for i, v := range vehicles {
		modelVehicles[i] = vehicleToModel(v)
	}

	return model.Problem{
		Depot:       depot,
		Nodes:       nodes,
		Vehicles:    modelVehicles,
		Durations:   durations,
		Distances:   distances,
		TimeLimitMS: 4500,
	}
}

func orderToNode(o models.Order) model.Node {
	nodeType := model.NodeTypeDelivery
	if o.Type == 1 {
		nodeType = model.NodeTypePickup
	}

	tags := make([]string, 0, len(o.Skills))
	for _, s := range o.Skills {
		tags = append(tags, s.Name)
	}

	twStart, twEnd := 0, 1440
	if o.TimeWindowStart != nil {
		twStart = *o.TimeWindowStart
	}
	if o.TimeWindowEnd != nil {
		twEnd = *o.TimeWindowEnd
	}

	serviceTime := 10
	if o.ServiceTime != nil {
		serviceTime = *o.ServiceTime
	}

	lat, lng := 0.0, 0.0
	if o.DesLatitude != nil {
		lat = *o.DesLatitude
	}
	if o.DesLongitude != nil {
		lng = *o.DesLongitude
	}

	return model.Node{
		ID:          strconv.FormatUint(uint64(o.ID), 10),
		Lat:         lat,
		Lng:         lng,
		Demand:      int(o.Capacity),
		ServiceTime: serviceTime,
		TWStart:     twStart,
		TWEnd:       twEnd,
		Tags:        tags,
		Type:        nodeType,
		Priority:    priorityFromInt(o.Priority),
	}
}

func vehicleToModel(v models.Vehicle) model.Vehicle {
	tags := make([]string, 0, len(v.Skills))
	for _, s := range v.Skills {
		tags = append(tags, s.Name)
	}

	shiftStart, shiftEnd := 480, 1020
	if v.DailyWorkTimeStart != nil {
		shiftStart = *v.DailyWorkTimeStart
	}
	if v.DailyWorkTimeEnd != nil {
		shiftEnd = *v.DailyWorkTimeEnd
	}

	breakStart, breakEnd := 0, 0
	if v.DailyBreakTimeStart != nil {
		breakStart = *v.DailyBreakTimeStart
	}
	if v.DailyBreakTimeEnd != nil {
		breakEnd = *v.DailyBreakTimeEnd
	}

	maxTasks := 20
	if v.MaxTask != nil {
		maxTasks = *v.MaxTask
	}

	startLat, startLng := 0.0, 0.0
	if v.StartLat != nil {
		startLat = *v.StartLat
	}
	if v.StartLon != nil {
		startLng = *v.StartLon
	}

	return model.Vehicle{
		ID:         strconv.FormatUint(uint64(v.ID), 10),
		Capacity:   int(v.Capacity),
		ShiftStart: shiftStart,
		ShiftEnd:   shiftEnd,
		BreakStart: breakStart,
		BreakEnd:   breakEnd,
		MaxTasks:   maxTasks,
		Tags:       tags,
		StartLat:   startLat,
		StartLng:   startLng,
		EndLat:     startLat,
		EndLng:     startLng,
	}
}

func buildRoutes(planID string, orders []models.Order, vehicles []models.Vehicle, solution model.Solution) []models.Route {
	vehicleByStrID := make(map[string]models.Vehicle, len(vehicles))
	for _, v := range vehicles {
		vehicleByStrID[strconv.FormatUint(uint64(v.ID), 10)] = v
	}

	orderByStrID := make(map[string]models.Order, len(orders))
	for _, o := range orders {
		orderByStrID[strconv.FormatUint(uint64(o.ID), 10)] = o
	}

	routes := make([]models.Route, 0, len(solution.Routes))
	for _, r := range solution.Routes {
		v, ok := vehicleByStrID[r.VehicleID]
		if !ok {
			continue
		}

		stops := make([]models.Stop, 0, len(r.Stops))
		for seq, s := range r.Stops {
			order, ok := orderByStrID[s.NodeID]
			if !ok {
				continue
			}
			stops = append(stops, models.Stop{
				SequenceNumber: seq + 1,
				OrderID:        order.ID,
				ArrivalMin:     s.ArrivalMin,
				DepartMin:      s.DepartMin,
			})
		}

		dist := r.TotalDistance
		dur := r.TotalDuration
		routes = append(routes, models.Route{
			VehicleID:     v.ID,
			PlanID:        planID,
			TotalDistance: &dist,
			TotalTime:     &dur,
			Stops:         stops,
		})
	}
	return routes
}

func intsToFloats(m [][]int) [][]float64 {
	result := make([][]float64, len(m))
	for i, row := range m {
		result[i] = make([]float64, len(row))
		for j, v := range row {
			result[i][j] = float64(v)
		}
	}
	return result
}

func priorityFromInt(p int) model.Priority {
	switch p {
	case 3:
		return model.PriorityCritical
	case 2:
		return model.PriorityHigh
	case 1:
		return model.PriorityMedium
	default:
		return model.PriorityLow
	}
}
