package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	dto "ROP_Backend/internal/dto/request"
	"ROP_Backend/internal/dto/response"

	"github.com/ROP-TEAM/rop-algorithm/model"
	"github.com/ROP-TEAM/rop-algorithm/solver"
)

type matrixBuilder interface {
	BuildMatrix(ctx context.Context, locs []dto.LocationInput) (response.MatrixResponse, error)
}

type PlanningService struct {
	matrixSvc matrixBuilder
	solver    solver.Solver
}

func NewPlanningService(
	matrixSvc matrixBuilder,
	solver solver.Solver,
) *PlanningService {
	return &PlanningService{
		matrixSvc: matrixSvc,
		solver:    solver,
	}
}

func validateRequest(req dto.OptimizeRequest) error {
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
	req dto.OptimizeRequest,
) (
	*response.OptimizeResponse,
	error,
) {
	if err := validateRequest(req); err != nil {
		return nil, err
	}

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
	), nil
}

func buildOptimizeLocations(
	req dto.OptimizeRequest,
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
	req dto.OptimizeRequest,
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

		TimeLimitMS: 4500,
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

	return model.Node{
		ID: strconv.Itoa(i),

		Lat: o.DesLatitude,

		Lng: o.DesLongitude,

		Demand: o.Capacity,

		ServiceTime: o.ServiceTime,

		TWStart: o.TimeWindowStart,

		TWEnd: o.TimeWindowEnd,

		Tags: o.Skills,

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

	return model.Vehicle{
		ID: strconv.Itoa(i),

		Capacity: v.Capacity,

		ShiftStart: v.DailyWorkTimeStart,

		ShiftEnd: v.DailyWorkTimeEnd,

		BreakStart: v.DailyBreakTimeStart,

		BreakEnd: v.DailyBreakTimeEnd,

		MaxTasks: v.MaxTask,

		Tags: v.Skills,

		StartLat: v.StartLatitude,

		StartLng: v.StartLongitude,

		EndLat: v.EndLatitude,

		EndLng: v.EndLongitude,
	}
}

func buildOptimizeResponse(
	req dto.OptimizeRequest,
	solution model.Solution,
) *response.OptimizeResponse {

	vehicleMap :=
		map[string]string{}

	orderMap :=
		map[string]string{}

	for i, v := range req.Vehicles {
		vehicleMap[strconv.Itoa(i)] = v.Name
	}

	for i, o := range req.Orders {
		orderMap[strconv.Itoa(i)] = o.Name
	}

	routes :=
		[]response.RouteResponse{}

	for _, r := range solution.Routes {

		stops :=
			[]response.StopResponse{}

		for _, stop := range r.Stops {

			stops = append(
				stops,
				response.StopResponse{
					OrderName: orderMap[stop.NodeID],

					ArrivalMin: stop.ArrivalMin,

					DepartMin: stop.DepartMin,
				},
			)
		}

		routes = append(
			routes,
			response.RouteResponse{
				VehicleName: vehicleMap[r.VehicleID],

				TotalDistance: r.TotalDistance,

				TotalDuration: float64(r.TotalDuration),

				Stops: stops,
			},
		)
	}

	unassigned := make([]string, 0, len(solution.Unassigned))
	for _, id := range solution.Unassigned {
		if name, ok := orderMap[id]; ok {
			unassigned = append(unassigned, name)
		}
	}

	dropReasons := make([]response.DropReasonResponse, 0, len(solution.DropReasons))
	for _, dr := range solution.DropReasons {
		dropReasons = append(dropReasons, response.DropReasonResponse{
			OrderName: orderMap[dr.NodeID],
			Code:      dr.Code,
			Detail:    dr.Detail,
		})
	}

	return &response.OptimizeResponse{
		Status:      string(solution.Status),
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
	case 4:
		return model.PriorityCritical
	case 3:
		return model.PriorityHigh
	case 2:
		return model.PriorityMedium
	default:
		return model.PriorityLow
	}
}
