package response

type MatrixTestResponse struct {
	Durations [][]int `json:"durations"`
	Distances [][]int `json:"distances"`
}
