package services

import (
	"context"
	"fmt"

	"ROP_Backend/internal/config"
	dto "ROP_Backend/internal/dto/request"
	"ROP_Backend/internal/dto/response"

	"github.com/ROP-TEAM/rop-algorithm/gmap"
	"github.com/ROP-TEAM/rop-algorithm/model"
	"github.com/ROP-TEAM/rop-algorithm/osrm"
)

type MatrixService struct {
	dm gmap.DistanceMatrix
}

func NewMatrixService(cfg *config.Config) (*MatrixService, error) {
	dm, err := newDistanceMatrix(cfg)
	if err != nil {
		return nil, err
	}
	return &MatrixService{dm: dm}, nil
}

func (s *MatrixService) BuildMatrix(locs []dto.LocationInput) (response.MatrixTestResponse, error) {
	modelLocs := make([]model.Location, len(locs))
	for i, l := range locs {
		modelLocs[i] = model.NewLatLngLocation(l.Lat, l.Lng)
	}
	durations, distances, err := s.dm.BuildMatrix(context.Background(), modelLocs, model.MatrixOptions{})
	if err != nil {
		return response.MatrixTestResponse{}, err
	}
	return response.MatrixTestResponse{Durations: durations, Distances: distances}, nil
}

func newDistanceMatrix(cfg *config.Config) (gmap.DistanceMatrix, error) {
	switch cfg.DISTANCE_MATRIX_PROVIDER {
	case "osrm":
		return osrm.NewOSRMMatrix(cfg.OSRM_BASE_URL)
	case "google", "":
		return gmap.NewGoogleMapsMatrix(cfg.GOOGLE_MAPS_API_KEY)
	default:
		return nil, fmt.Errorf("unknown DISTANCE_MATRIX_PROVIDER: %s", cfg.DISTANCE_MATRIX_PROVIDER)
	}
}
