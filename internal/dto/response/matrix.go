package response

type MatrixResponse struct {
	Durations [][]int `json:"durations"`
	Distances [][]int `json:"distances"`
}

type BuildMatrixResponse struct {
	Message string         `json:"message"`
	Node    int            `json:"node"`
	Result  MatrixResponse `json:"result"`
}
