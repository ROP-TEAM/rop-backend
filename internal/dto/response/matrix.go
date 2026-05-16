package response

type MatrixResponse struct {
	Durations [][]int `json:"durations"`
	Distances [][]int `json:"distances"`
}
