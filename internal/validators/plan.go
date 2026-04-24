package validators

import (
	"ROP_Backend/internal/dto"
	"ROP_Backend/internal/utils"
	"fmt"
)

func ValidatePlan(req dto.PlanRequest) error {
	vMap := make(map[int]dto.VehiclesRequest)

	for _, v := range req.Vehicles {
		if v.Capacity <= 0 {
			return fmt.Errorf("vehicle %d has non-positive capacity", v.ID)
		}
		if v.MaxTask <= 0 {
			return fmt.Errorf("vehicle %d has non-positive maxTask", v.ID)
		}
		if v.WorkTime.Start >= v.WorkTime.End {
			return fmt.Errorf("vehicle %d has invalid work time", v.ID)
		}
		if v.BreakTime.Start >= v.BreakTime.End {
			return fmt.Errorf("vehicle %d has invalid break time", v.ID)
		}

		vMap[v.ID] = v
	}

	for _, o := range req.Orders {
		if o.Capacity <= 0 {
			return fmt.Errorf("order %d has non-positive capacity", o.ID)
		}
		if o.TimeWindow.Start >= o.TimeWindow.End {
			return fmt.Errorf("order %d has invalid time window", o.ID)
		}
		if o.DeliveryWindow.Start >= o.DeliveryWindow.End {
			return fmt.Errorf("order %d has invalid delivery window", o.ID)
		}
		if o.ServiceTime < 0 {
			return fmt.Errorf("order %d has negative service time", o.ID)
		}
		if o.Type == "" {
			return fmt.Errorf("order %d has empty type", o.ID)
		}

		// Check vehicle assignment
		v, ok := vMap[o.AssignToVehicleID]
		if !ok {
			return fmt.Errorf("order %d invalid vehicle assign", o.ID)
		}
		if !utils.IsSkillMatch(o.Skills, v.Skills) {
			return fmt.Errorf("order %d skill mismatch", o.ID)
		}
		if o.Capacity > v.Capacity {
			return fmt.Errorf("order %d exceed capacity", o.ID)
		}
		if !utils.IsTimeOverlap(o.TimeWindow, v.WorkTime) {
			return fmt.Errorf("order %d outside workTime", o.ID)
		}
		if utils.IsTimeOverlap(o.TimeWindow, v.BreakTime) {
			return fmt.Errorf("order %d overlaps breakTime", o.ID)
		}
	}
	return nil
}
