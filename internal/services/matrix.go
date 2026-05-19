package services

import (
	"context"
	"fmt"

	"ROP_Backend/internal/config"
	dto "ROP_Backend/internal/dto/request"
	"ROP_Backend/internal/dto/response"
	"ROP_Backend/internal/models"

	"ROP_Backend/internal/providers/gmap"
	"ROP_Backend/internal/providers/osrm"
)

type MatrixService struct {
	dm       gmap.DistanceMatrix
	provider string
}

func NewMatrixService(cfg *config.Config) (*MatrixService, error) {
	dm, err := newDistanceMatrix(cfg)
	if err != nil {
		return nil, err
	}
	provider := cfg.DISTANCE_MATRIX_PROVIDER
	if provider == "" {
		provider = "google"
	}
	return &MatrixService{dm: dm, provider: provider}, nil
}

func (s *MatrixService) Provider() string {
	return s.provider
}

func (s *MatrixService) BuildMatrix(ctx context.Context, locs []dto.LocationInput) (response.MatrixResponse, error) {
	modelLocs := make([]models.Location, len(locs))
	for i, l := range locs {
		modelLocs[i] = models.NewLatLngLocation(l.Lat, l.Lng)
	}
	durations, distances, err := s.dm.BuildMatrix(ctx, modelLocs, models.MatrixOptions{})
	if err != nil {
		return response.MatrixResponse{}, err
	}
	return response.MatrixResponse{Durations: durations, Distances: distances}, nil
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
