package service

import (
	"github.com/raffops/gofipe/cmd/goFipe/domain"
	"github.com/raffops/gofipe/cmd/goFipe/errs"
	"github.com/raffops/gofipe/cmd/goFipe/port"
)

type VehicleService struct {
	vehicleRepo port.VehicleRepository
}

func NewVehicleService(vehicleRepo port.VehicleRepository) VehicleService {
	return VehicleService{vehicleRepo: vehicleRepo}
}

func (v VehicleService) GetVehicleByFipeCode(fipeCode int,
	orderBy []domain.OrderBy,
	pagination domain.Pagination) ([]domain.Vehicle, *errs.AppError) {
	conditions := []domain.Condition{{Column: "fipe_code", Value: fipeCode, Operator: "="}}
	return v.vehicleRepo.GetVehicle(conditions, orderBy, pagination)
}

func (v VehicleService) GetVehicleByReferenceYearMonth(
	year int,
	month int,
	orderBy []domain.OrderBy,
	pagination domain.Pagination) ([]domain.Vehicle, *errs.AppError) {

	conditions := []domain.Condition{
		{Column: "year", Value: year, Operator: "="},
		{Column: "month", Value: month, Operator: "="},
	}

	return v.vehicleRepo.GetVehicle(conditions, orderBy, pagination)
}
