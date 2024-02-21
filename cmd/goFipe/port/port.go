package port

import (
	"github.com/raffops/gofipe/cmd/goFipe/domain"
	"github.com/raffops/gofipe/cmd/goFipe/errs"
)

//go:generate mockgen -source port.go -destination ../mocks/mockVehicle.go

type VehicleService interface {
	GetVehicle(where map[string]string, orderBy map[string]bool, limit int, offset int) ([]domain.Vehicle, *errs.AppError)
}

type VehicleRepository interface {
	GetVehicle(conditions []domain.WhereClause, orderBy []domain.OrderByClause, pagination domain.Pagination) ([]domain.Vehicle, *errs.AppError)
}
