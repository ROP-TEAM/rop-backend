package services

import (
	"fmt"

	"github.com/ROP-TEAM/rop-algorithm/gmap"
	"github.com/ROP-TEAM/rop-algorithm/osrm"
	"ROP_Backend/internal/config"
)

// NewDistanceMatrix creates the right DistanceMatrix implementation based on DISTANCE_MATRIX_PROVIDER.
// Supported values: "google" (default), "osrm".
func NewDistanceMatrix(cfg *config.Config) (gmap.DistanceMatrix, error) {
	switch cfg.DISTANCE_MATRIX_PROVIDER {
	case "osrm":
		return osrm.NewOSRMMatrix(cfg.OSRM_BASE_URL)
	case "google", "":
		return gmap.NewGoogleMapsMatrix(cfg.GOOGLE_MAPS_API_KEY)
	default:
		return nil, fmt.Errorf("unknown DISTANCE_MATRIX_PROVIDER: %s", cfg.DISTANCE_MATRIX_PROVIDER)
	}
}
