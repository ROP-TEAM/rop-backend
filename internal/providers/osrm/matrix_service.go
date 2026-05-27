package osrm

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"ROP_Backend/internal/models"
)

// OSRMMatrix implements gmap.DistanceMatrix via the OSRM /table API.
type OSRMMatrix struct {
	baseURL    string
	httpClient *http.Client
	chunkSize  int
}

// OSRMMatrixOption configures an OSRMMatrix.
type OSRMMatrixOption func(*OSRMMatrix)

// WithHTTPClient sets the HTTP client used to call OSRM.
func WithHTTPClient(client *http.Client) OSRMMatrixOption {
	return func(m *OSRMMatrix) {
		m.httpClient = client
	}
}

// WithChunkSize overrides the row chunk size for parallel BuildMatrix requests (default 100).
func WithChunkSize(size int) OSRMMatrixOption {
	return func(m *OSRMMatrix) {
		m.chunkSize = size
	}
}

// NewOSRMMatrix creates an OSRMMatrix. baseURL is the OSRM server root, e.g. "http://localhost:5000".
func NewOSRMMatrix(baseURL string, opts ...OSRMMatrixOption) (*OSRMMatrix, error) {
	if strings.TrimSpace(baseURL) == "" {
		return nil, fmt.Errorf("osrm base URL is required")
	}
	m := &OSRMMatrix{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(m)
		}
	}
	return m, nil
}

// BuildMatrix returns n×n matrices of durations (minutes) and distances (meters).
// Requests with more rows than chunkSize are split into parallel chunk requests.
func (m *OSRMMatrix) BuildMatrix(ctx context.Context, locs []models.Location, opts models.MatrixOptions) ([][]int, [][]int, error) {
	n := len(locs)
	if n == 0 {
		return [][]int{}, [][]int{}, nil
	}

	chunkSize := m.chunkSize
	if chunkSize <= 0 {
		chunkSize = defaultChunkSize
	}

	if n <= chunkSize {
		return m.buildMatrixSingle(ctx, locs, n)
	}
	return m.buildMatrixParallel(ctx, locs, chunkSize)
}

func (m *OSRMMatrix) buildMatrixSingle(ctx context.Context, locs []models.Location, n int) ([][]int, [][]int, error) {
	url := fmt.Sprintf("%s/table/v1/driving/%s?annotations=duration,distance", m.baseURL, formatCoordinates(locs))

	resp, err := m.executeHTTP(ctx, url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := readAndCheckHTTP(resp)
	if err != nil {
		return nil, nil, err
	}

	return parseTableResponse(body, n)
}
