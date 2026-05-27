package osrm

import (
	"ROP_Backend/internal/models"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// full 3×3 reference data (seconds / meters)
var refDurations = [][]float64{
	{0, 120, 240},
	{120, 0, 180},
	{240, 180, 0},
}
var refDistances = [][]float64{
	{0, 1200, 2400},
	{1200, 0, 1800},
	{2400, 1800, 0},
}

func parseSourceParams(raw string) []int {
	// r.URL.Query() already URL-decodes %3B → ";" before returning
	parts := strings.Split(raw, ";")
	var indices []int
	for _, p := range parts {
		i, err := strconv.Atoi(strings.TrimSpace(p))
		if err == nil {
			indices = append(indices, i)
		}
	}
	return indices
}

func chunkMatrixHandler(w http.ResponseWriter, r *http.Request) {
	indices := parseSourceParams(r.URL.Query().Get("sources"))

	dur := make([][]float64, len(indices))
	dist := make([][]float64, len(indices))
	for i, idx := range indices {
		dur[i] = refDurations[idx]
		dist[i] = refDistances[idx]
	}

	json.NewEncoder(w).Encode(OSRMTableResponse{Code: "Ok", Durations: dur, Distances: dist})
}

func TestBuildMatrixChunked(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(chunkMatrixHandler))
	defer server.Close()

	m, err := NewOSRMMatrix(server.URL, WithChunkSize(2))
	if err != nil {
		t.Fatalf("NewOSRMMatrix: %v", err)
	}

	locs := []models.Location{
		models.NewLatLngLocation(13.75, 100.50),
		models.NewLatLngLocation(13.74, 100.53),
		models.NewLatLngLocation(13.76, 100.55),
	}

	durations, distances, err := m.BuildMatrix(context.Background(), locs, models.MatrixOptions{})
	if err != nil {
		t.Fatalf("BuildMatrix: %v", err)
	}

	if len(durations) != 3 || len(distances) != 3 {
		t.Fatalf("expected 3×3 matrices, got %d×%d", len(durations), len(distances))
	}

	// seconds→minutes (floor): 120→2, 240→4, 180→3
	wantDur := [][]int{{0, 2, 4}, {2, 0, 3}, {4, 3, 0}}
	wantDist := [][]int{{0, 1200, 2400}, {1200, 0, 1800}, {2400, 1800, 0}}

	for i := range 3 {
		for j := range 3 {
			if durations[i][j] != wantDur[i][j] {
				t.Errorf("durations[%d][%d] = %d, want %d", i, j, durations[i][j], wantDur[i][j])
			}
			if distances[i][j] != wantDist[i][j] {
				t.Errorf("distances[%d][%d] = %d, want %d", i, j, distances[i][j], wantDist[i][j])
			}
		}
	}
}

func TestBuildMatrixChunkErrorPropagates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		indices := parseSourceParams(r.URL.Query().Get("sources"))
		for _, idx := range indices {
			if idx >= 2 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		json.NewEncoder(w).Encode(OSRMTableResponse{
			Code:      "Ok",
			Durations: [][]float64{{0, 120}, {120, 0}},
			Distances: [][]float64{{0, 1200}, {1200, 0}},
		})
	}))
	defer server.Close()

	m, err := NewOSRMMatrix(server.URL, WithChunkSize(2))
	if err != nil {
		t.Fatalf("NewOSRMMatrix: %v", err)
	}

	locs := []models.Location{
		models.NewLatLngLocation(13.75, 100.50),
		models.NewLatLngLocation(13.74, 100.53),
		models.NewLatLngLocation(13.76, 100.55),
	}

	_, _, err = m.BuildMatrix(context.Background(), locs, models.MatrixOptions{})
	if err == nil {
		t.Fatal("expected error when a chunk fails")
	}
}
