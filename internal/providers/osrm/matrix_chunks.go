package osrm

import (
	"context"
	"fmt"
	"sync"

	"ROP_Backend/internal/models"
	"golang.org/x/sync/errgroup"
)

const defaultChunkSize = 100

func (m *OSRMMatrix) buildMatrixParallel(ctx context.Context, locs []models.Location, chunkSize int) ([][]int, [][]int, error) {
	n := len(locs)
	allCoords := formatCoordinates(locs)
	allDests := indexRange(0, n)

	durations := make([][]int, n)
	distances := make([][]int, n)
	for i := range n {
		durations[i] = make([]int, n)
		distances[i] = make([]int, n)
	}

	var mu sync.Mutex
	g, gctx := errgroup.WithContext(ctx)

	for start := 0; start < n; start += chunkSize {
		chunkStart := start
		chunkEnd := min(start+chunkSize, n)
		rowCount := chunkEnd - chunkStart

		g.Go(func() error {
			url := fmt.Sprintf("%s/table/v1/driving/%s?annotations=duration,distance&sources=%s&destinations=%s",
				m.baseURL, allCoords, indexRange(chunkStart, chunkEnd), allDests)

			resp, err := m.executeHTTP(gctx, url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			body, err := readAndCheckHTTP(resp)
			if err != nil {
				return err
			}

			chunkDur, chunkDist, err := parseTableResponseWithShape(body, rowCount, n)
			if err != nil {
				return err
			}

			mu.Lock()
			for i, row := range chunkDur {
				copy(durations[chunkStart+i], row)
			}
			for i, row := range chunkDist {
				copy(distances[chunkStart+i], row)
			}
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, nil, err
	}
	return durations, distances, nil
}
