package response

// WastePlanResponse is the week plan scored against the current assignment.
type WastePlanResponse struct {
	Baseline WasteBaseline    `json:"baseline"`
	PerDay   []WasteDayMetric `json:"perDay"`
	Totals   WasteTotals      `json:"totals"`
}

// WasteBaseline is the current vehicle assignment scored with the same method as
// the plan, so the deadhead is directly comparable.
type WasteBaseline struct {
	TruckGroups        int     `json:"truckGroups"`
	AssignedUnits      int     `json:"assignedUnits"`
	InterStopDeadheadKm float64 `json:"interStopDeadheadKm"`
}

// WasteDayMetric is one weekday's planned outcome.
type WasteDayMetric struct {
	Weekday            int     `json:"weekday"`
	Cells              int     `json:"cells"`
	Bands              int     `json:"bands"`
	TrucksUsed         int     `json:"trucksUsed"`
	AssignedUnits      int     `json:"assignedUnits"`
	UnassignedUnits    int     `json:"unassignedUnits"`
	InterStopDeadheadKm float64 `json:"interStopDeadheadKm"`
	MeanUtilisation    float64 `json:"meanUtilisation"`
	LoadGini           float64 `json:"loadGini"`
}

// WasteTotals aggregates the week and its improvement over the baseline.
type WasteTotals struct {
	TruckShifts           int     `json:"truckShifts"`
	AssignedUnits         int     `json:"assignedUnits"`
	UnassignedUnits       int     `json:"unassignedUnits"`
	InterStopDeadheadKm   float64 `json:"interStopDeadheadKm"`
	InterStopReductionPct float64 `json:"interStopReductionPct"`
}
