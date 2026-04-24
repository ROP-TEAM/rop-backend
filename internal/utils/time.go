package utils

import "ROP_Backend/internal/dto"

func IsTimeOverlap(a, b dto.TimeRange) bool {
	return a.Start < b.End && b.Start < a.End
}
