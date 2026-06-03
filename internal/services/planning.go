package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	dto "ROP_Backend/internal/dto/request"
	"ROP_Backend/internal/dto/response"
	"ROP_Backend/internal/models"
	"ROP_Backend/internal/repository"

	"github.com/ROP-TEAM/rop-algorithm/model"
	"github.com/ROP-TEAM/rop-algorithm/solver"
)

type matrixBuilder interface {
	BuildMatrix(ctx context.Context, locs []dto.LocationInput) (response.MatrixResponse, error)
}

type skillInfo struct {
	ID    int
	Color string
}

type PlanningService struct {
	matrixSvc matrixBuilder
	solver    solver.Solver

	planRepo    *repository.PlanRepository
	vehicleRepo *repository.VehicleRepository
	orderRepo   *repository.OrderRepository
	tagRepo     *repository.TagSkillRepository
}

func NewPlanningService(
	matrixSvc matrixBuilder,
	solver solver.Solver,
	planRepo *repository.PlanRepository,
	vehicleRepo *repository.VehicleRepository,
	orderRepo *repository.OrderRepository,
	tagRepo *repository.TagSkillRepository,
) *PlanningService {
	return &PlanningService{
		matrixSvc:   matrixSvc,
		solver:      solver,
		planRepo:    planRepo,
		vehicleRepo: vehicleRepo,
		orderRepo:   orderRepo,
		tagRepo:     tagRepo,
	}
}

func validateRequest(req *dto.OptimizeRequest) error {
	if len(req.Vehicles) == 0 {
		return errors.New("vehicles cannot be empty")
	}
	if len(req.Orders) == 0 {
		return errors.New("orders cannot be empty")
	}
	if req.DepotLat < -90 || req.DepotLat > 90 {
		return errors.New("depot_lat must be between -90 and 90")
	}
	if req.DepotLon < -180 || req.DepotLon > 180 {
		return errors.New("depot_lon must be between -180 and 180")
	}
	for i, o := range req.Orders {
		if o.DesLatitude < -90 || o.DesLatitude > 90 {
			return fmt.Errorf("orders[%d]: des_latitude must be between -90 and 90", i)
		}
		if o.DesLongitude < -180 || o.DesLongitude > 180 {
			return fmt.Errorf("orders[%d]: des_longitude must be between -180 and 180", i)
		}
		if o.TimeWindowStart > o.TimeWindowEnd {
			return fmt.Errorf("orders[%d]: time_window_start must not exceed time_window_end", i)
		}
		if o.Capacity < 0 {
			return fmt.Errorf("orders[%d]: capacity must be >= 0", i)
		}
		if o.ServiceTime < 0 {
			return fmt.Errorf("orders[%d]: service_time must be >= 0", i)
		}
	}
	return nil
}

func (s *PlanningService) Optimize(
	ctx context.Context,
	req *dto.OptimizeRequest,
) (
	*response.OptimizeResponse,
	error,
) {
	if err := validateRequest(req); err != nil {
		return nil, err
	}

	// begin persis to company
	plan := &models.Plan{
		CompanyID: "c5c5cfc5-d97e-4ade-ae2f-7abec89c6f12",
	}

	if err := s.planRepo.Create(ctx, plan); err != nil {
		return nil, err
	}

	planID := plan.ID

	skillMap := map[string]skillInfo{}

	for _, v := range req.Vehicles {
		for _, skill := range v.Skills {
			if _, ok := skillMap[skill.Name]; ok {
				continue
			}
			color := skill.Color
			if color == "" {
				color = "#3B82F6"
			}

			tagSkill := models.TagSkill{
				Name:   skill.Name,
				Color:  color,
				PlanID: planID,
			}

			if err := s.tagRepo.Create(&tagSkill); err != nil {
				return nil, err
			}

			skillMap[skill.Name] = skillInfo{
				ID:    int(tagSkill.ID),
				Color: color,
			}
		}
	}

	for _, o := range req.Orders {
		if o.Skill == "" {
			continue
		}

		if _, ok := skillMap[o.Skill]; ok {
			continue
		}

		tagSkill := models.TagSkill{
			Name:   o.Skill,
			Color:  "#3B82F6",
			PlanID: planID,
		}

		if err := s.tagRepo.Create(&tagSkill); err != nil {
			return nil, err
		}

		skillMap[o.Skill] = skillInfo{
			ID:    int(tagSkill.ID),
			Color: "#3B82F6",
		}
	}

	vehicles := make([]models.Vehicle, len(req.Vehicles))

	for i, v := range req.Vehicles {

		// ทำไปทำไมนะ
		// skills := make([]models.TagSkill, 0, len(v.Skills))

		// for _, s := range v.Skills {
		// 	skills = append(skills, models.TagSkill{
		// 		ID: uint(skillMap[s.Name].ID),
		// 	})
		// }

		vehicles[i] = models.Vehicle{
			Name:        v.Name,
			Model:       v.Model,
			PlateNumber: v.PlateNumber,

			Capacity: float64(v.Capacity),
			MaxTask:  &v.MaxTask,

			PlanID: planID,

			// StartLat: &v.StartLatitude,
			// StartLon: &v.StartLongitude,
			// EndLat:   &v.EndLatitude,
			// EndLon:   &v.EndLongitude,

			DailyWorkTimeStart:  &v.DailyWorkTimeStart,
			DailyWorkTimeEnd:    &v.DailyWorkTimeEnd,
			DailyBreakTimeStart: &v.DailyBreakTimeStart,
			DailyBreakTimeEnd:   &v.DailyBreakTimeEnd,
		}
	}

	if err := s.vehicleRepo.BatchCreate(ctx, vehicles); err != nil {
		return nil, err
	}

	orders := make([]models.Order, len(req.Orders))

	for i, o := range req.Orders {

		orders[i] = models.Order{
			Name: o.Name,

			Type: o.Type,

			Capacity: float64(o.Capacity),

			Priority: o.Priority,

			PlanID: planID,

			DesLatitude:  &o.DesLatitude,
			DesLongitude: &o.DesLongitude,

			ServiceTime: &o.ServiceTime,

			TimeWindowStart: &o.TimeWindowStart,
			TimeWindowEnd:   &o.TimeWindowEnd,
		}
	}

	if err := s.orderRepo.BatchCreate(ctx, orders); err != nil {
		return nil, err
	}
	// end of persis to company

	locations := buildOptimizeLocations(req)

	matrixResp, err := s.matrixSvc.BuildMatrix(ctx, locations)

	if err != nil {
		return nil, err
	}

	problem := buildOptimizeProblem(
		req,
		intsToFloats(
			matrixResp.Durations,
		),
		intsToFloats(
			matrixResp.Distances,
		),
	)

	solution, err :=
		s.solver.Solve(
			ctx,
			problem,
		)

	if err != nil {
		return nil, err
	}

	return buildOptimizeResponse(
		req,
		solution,
		matrixResp.Distances,
		matrixResp.Durations,
		orders,
		vehicles,
		skillMap,
	), nil
}

func buildOptimizeLocations(
	req *dto.OptimizeRequest,
) []dto.LocationInput {

	locs := []dto.LocationInput{
		{
			Lat: req.DepotLat,
			Lng: req.DepotLon,
		},
	}

	for _, o := range req.Orders {
		locs = append(
			locs,
			dto.LocationInput{
				Lat: o.DesLatitude,
				Lng: o.DesLongitude,
			},
		)
	}

	return locs
}

func buildOptimizeProblem(
	req *dto.OptimizeRequest,
	durations [][]float64,
	distances [][]float64,
) model.Problem {

	nodes := make(
		[]model.Node,
		len(req.Orders),
	)

	for i, o := range req.Orders {
		nodes[i] =
			optimizeOrderToNode(
				i,
				o,
			)
	}

	vehicles := make(
		[]model.Vehicle,
		len(req.Vehicles),
	)

	for i, v := range req.Vehicles {
		vehicles[i] =
			optimizeVehicleToModel(
				i,
				v,
			)
	}

	return model.Problem{
		Depot: model.Node{
			ID: "depot",

			Lat: req.DepotLat,
			Lng: req.DepotLon,

			Type: model.NodeTypeDepot,
		},

		Nodes: nodes,

		Vehicles: vehicles,

		Durations: durations,

		Distances: distances,

		TimeLimitMS:     req.TimeLimitMS,
		EnableALNS:      req.EnableALNS,
		EnableMultiTrip: req.EnableMultiTrip,
		ReloadMin:       req.ReloadMin,
		Seed:            req.Seed,
	}
}

func optimizeOrderToNode(
	i int,
	o dto.OptimizeOrder,
) model.Node {

	nodeType :=
		model.NodeTypeDelivery

	if o.Type == 1 {
		nodeType =
			model.NodeTypePickup
	}

	tags := []string{}

	if o.Skill != "" {
		tags = append(tags, o.Skill)
	}

	return model.Node{
		ID: strconv.Itoa(i),

		Lat: o.DesLatitude,

		Lng: o.DesLongitude,

		Demand: o.Capacity,

		ServiceTime: o.ServiceTime,

		TWStart: o.TimeWindowStart,

		TWEnd: o.TimeWindowEnd,

		Tags: tags,

		Type: nodeType,

		Priority: priorityFromInt(
			o.Priority,
		),
	}
}

func optimizeVehicleToModel(
	i int,
	v dto.OptimizeVehicle,
) model.Vehicle {
	tags := make([]string, 0, len(v.Skills))

	for _, s := range v.Skills {
		tags = append(tags, s.Name)
	}

	return model.Vehicle{
		ID: strconv.Itoa(i),

		Capacity: v.Capacity,

		ShiftStart: v.DailyWorkTimeStart,

		ShiftEnd: v.DailyWorkTimeEnd,

		BreakStart: v.DailyBreakTimeStart,

		BreakEnd: v.DailyBreakTimeEnd,

		MaxTasks: v.MaxTask,

		Tags: tags,

		// StartLat: v.StartLatitude,

		// StartLng: v.StartLongitude,

		// EndLat: v.EndLatitude,

		// EndLng: v.EndLongitude,
	}
}

func buildOptimizeResponse(
	req *dto.OptimizeRequest,
	solution model.Solution,
	distances [][]int,
	durations [][]int,
	savedOrders []models.Order,
	savedVehicles []models.Vehicle,
	skillMap map[string]skillInfo,
) *response.OptimizeResponse {
	routes := make([]response.RouteResponse, 0, len(solution.Routes))

	for _, r := range solution.Routes {

		stops := make([]response.StopResponse, 0, len(r.Stops))
		prevMatrixIdx := 0 // depot

		for _, stop := range r.Stops {
			nodeIDInt, err := strconv.Atoi(stop.NodeID)

			if err != nil || nodeIDInt < 0 || nodeIDInt >= len(req.Orders) {
				continue
			}

			currMatrixIdx := nodeIDInt + 1

			var distFromPrev float64
			var timeFromPrev int
			if prevMatrixIdx < len(distances) && currMatrixIdx < len(distances[prevMatrixIdx]) {
				distFromPrev = float64(distances[prevMatrixIdx][currMatrixIdx])
			}
			if prevMatrixIdx < len(durations) && currMatrixIdx < len(durations[prevMatrixIdx]) {
				timeFromPrev = durations[prevMatrixIdx][currMatrixIdx]
			}

			order := req.Orders[nodeIDInt]
			dbOrder := savedOrders[nodeIDInt]

			var skillPtr *string
			if order.Skill != "" {
				skillPtr = &req.Orders[nodeIDInt].Skill
			}

			stops = append(stops, response.StopResponse{
				OrderName:            order.Name,
				ArrivalMin:           stop.ArrivalMin,
				DistanceFromPrevious: distFromPrev,
				DurationFromPrevious: timeFromPrev,

				// add from frontend requirement
				Capacity:        order.Capacity,
				TimeWindowStart: order.TimeWindowStart,
				TimeWindowEnd:   order.TimeWindowEnd,
				DesLatitude:     order.DesLatitude,
				DesLongitude:    order.DesLongitude,
				ServiceTime:     order.ServiceTime,
				Type:            order.Type,
				Priority:        order.Priority,
				Skill:           skillPtr,
				ID:              int(dbOrder.ID),
				Note:            nil,
			})

			prevMatrixIdx = currMatrixIdx
		}

		routeResp := response.RouteResponse{
			TotalDistance: r.TotalDistance,
			TotalDuration: float64(r.TotalDuration),
			Stops:         stops,
		}

		if vIdx, err := strconv.Atoi(r.VehicleID); err == nil && vIdx >= 0 && vIdx < len(req.Vehicles) {
			v := req.Vehicles[vIdx]
			dbVehicle := savedVehicles[vIdx]

			routeResp.ID = int(dbVehicle.ID)
			routeResp.Name = v.Name
			routeResp.Capacity = v.Capacity
			routeResp.WorkTimeStart = v.DailyWorkTimeStart
			routeResp.WorkTimeEnd = v.DailyWorkTimeEnd

			if v.Model != "" {
				routeResp.Model = &req.Vehicles[vIdx].Model
			}
			if v.PlateNumber != "" {
				routeResp.PlateNumber = &req.Vehicles[vIdx].PlateNumber
			}

			maxTaskCopy := v.MaxTask
			routeResp.MaxTask = &maxTaskCopy

			breakStartCopy := v.DailyBreakTimeStart
			routeResp.BreakTimeStart = &breakStartCopy

			breakEndCopy := v.DailyBreakTimeEnd
			routeResp.BreakTimeEnd = &breakEndCopy

			if len(v.Skills) > 0 {
				routeResp.Skills = make([]response.VehicleSkill, len(v.Skills))
				for idx, sk := range v.Skills {
					info := skillMap[sk.Name]
					idCopy := info.ID
					colorCopy := info.Color

					routeResp.Skills[idx] = response.VehicleSkill{
						Name:  sk.Name,
						ID:    &idCopy,
						Color: &colorCopy,
					}
				}
			}
		}

		routes = append(routes, routeResp)
	}

	unassigned := make([]string, 0, len(solution.Unassigned))
	for _, idStr := range solution.Unassigned {
		if idx, err := strconv.Atoi(idStr); err == nil && idx >= 0 && idx < len(req.Orders) {
			unassigned = append(unassigned, req.Orders[idx].Name)
		} else {
			unassigned = append(unassigned, idStr) // Fallback to raw ID string if unparsable
		}
	}

	dropReasons := make([]response.DropReasonResponse, 0, len(solution.DropReasons))
	for _, dr := range solution.DropReasons {
		dropReasons = append(dropReasons, response.DropReasonResponse{
			Code:   dr.Code,
			Detail: dr.Detail,
		})
	}

	return &response.OptimizeResponse{
		Message:     string(solution.Status),
		Routes:      routes,
		Unassigned:  unassigned,
		DropReasons: dropReasons,
	}
}

func intsToFloats(in [][]int) [][]float64 {
	out := make([][]float64, len(in))
	for i, row := range in {
		out[i] = make([]float64, len(row))
		for j, v := range row {
			out[i][j] = float64(v)
		}
	}
	return out
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
