package domain

//go:generate mockgen -source vehicle.go -destination ../mocks/domain/mockVehicle.go

type Vehicle struct {
	ReferenceMonth int
	FipeCode       int
	Brand          string
	Model          string
	YearModel      string
	Authentication string
	MeanValue      float64
}

type VehicleService interface {
	GetVehicleByFipeCode(int) ([]Vehicle, error)
	GetVehicleByReferenceMonth(int) ([]Vehicle, error)
}

type VehicleRepository interface {
	GetVehicleByFipeCode(int) ([]Vehicle, error)
	GetVehicleByReferenceMonth(int) ([]Vehicle, error)
}
