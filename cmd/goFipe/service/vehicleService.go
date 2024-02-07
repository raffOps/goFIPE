package service

import (
	"github.com/raffops/gofipe/cmd/goFipe/domain"
	"github.com/raffops/gofipe/cmd/goFipe/errs"
)

type VehicleService struct {
	vehicleRepo domain.VehicleRepository
}

func NewVehicleService(vehicleRepo domain.VehicleRepository) VehicleService {
	return VehicleService{vehicleRepo: vehicleRepo}
}

func (v VehicleService) GetVehicleByFipeCode(fipeCode int,
	orderBy []domain.OrderBy,
	pagination domain.Pagination) ([]domain.Vehicle, *errs.AppError) {
	conditions := []domain.Condition{{Column: "fipe_code", Value: fipeCode}}
	return v.vehicleRepo.GetVehicle(conditions, orderBy, pagination)
}

func (v VehicleService) GetVehicleByReferenceYearMonth(year int,
	month int,
	orderBy []domain.OrderBy,
	pagination domain.Pagination) ([]domain.Vehicle, *errs.AppError) {

	conditions := []domain.Condition{
		{Column: "year", Value: year},
		{Column: "month", Value: month},
	}

	return v.vehicleRepo.GetVehicle(conditions, orderBy, pagination)
}
