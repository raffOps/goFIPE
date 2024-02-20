package port

import (
	"github.com/raffops/gofipe/cmd/goFipe/domain"
	"github.com/raffops/gofipe/cmd/goFipe/errs"
)

type VehicleService interface {
	GetVehicleByFipeCode(fipeCode int, orderBy []domain.OrderBy, pagination domain.Pagination) ([]domain.Vehicle, *errs.AppError)
	GetVehicleByReferenceYearMonth(year int, month int, orderBy []domain.OrderBy, pagination domain.Pagination) ([]domain.Vehicle, *errs.AppError)
}

type VehicleRepository interface {
	GetVehicle(conditions []domain.Condition, orderBy []domain.OrderBy, pagination domain.Pagination) ([]domain.Vehicle, *errs.AppError)
}
