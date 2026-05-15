package dto

type MatrixTestRequest struct {
	Locations []LocationInput `json:"locations"`
}

type LocationInput struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
