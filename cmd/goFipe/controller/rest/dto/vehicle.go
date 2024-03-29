package dto

import "github.com/raffops/gofipe/cmd/goFipe/domain"

type GetVehicleResponse struct {
	Year           int     `json:"ano"`
	Month          int     `json:"mes"`
	FipeCode       string  `json:"fipe_code"`
	Brand          string  `json:"marca"`
	Model          string  `json:"modelo"`
	YearModel      string  `json:"ano_modelo"`
	Authentication string  `json:"autenticacao"`
	MeanValue      float32 `json:"valor_medio"`
}

func VehicleResponseFromDomain(vehicle domain.Vehicle) GetVehicleResponse {
	return GetVehicleResponse{
		Year:           vehicle.Year,
		Month:          vehicle.Month,
		FipeCode:       vehicle.FipeCode,
		Brand:          vehicle.Brand,
		Model:          vehicle.Model,
		YearModel:      vehicle.YearModel,
		Authentication: vehicle.Authentication,
		MeanValue:      vehicle.MeanValue,
	}
}
