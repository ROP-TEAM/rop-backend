package dto

import "github.com/ROP-TEAM/rop-algorithm/waste"

// WastePlanRequest plans a full week of garbage collection from cleaned route
// units. Each unit is one collection route treated as an atomic truck load
// (Phase 1). Zero-valued knobs fall back to the planner defaults.
type WastePlanRequest struct {
	Units []waste.RouteUnit `json:"units" validate:"required,min=1,dive"`

	CollectKmh float64 `json:"collectKmh"` // 0 → 5.0
	DriveKmh   float64 `json:"driveKmh"`   // 0 → 20.0

	ShiftStart int `json:"shiftStart" validate:"gte=0"`       // minutes from midnight
	ShiftEnd   int `json:"shiftEnd" validate:"gtefield=ShiftStart"` // minutes from midnight

	MaxPerCell  int `json:"maxPerCell"`  // 0 → 300
	MaxShiftMin int `json:"maxShiftMin"` // 0 → waste.DefaultMaxShiftMin
	TimeLimitMS int `json:"timeLimitMS"` // 0 → solver default

	Mode  string `json:"mode" validate:"omitempty,oneof=min full"` // default min
	Count int    `json:"count"`                                    // trucks per cell for full mode
}
