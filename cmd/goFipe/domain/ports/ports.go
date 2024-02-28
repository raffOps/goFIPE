package ports

import (
	"github.com/raffops/gofipe/cmd/goFipe/domain"
	"github.com/raffops/gofipe/cmd/goFipe/errs"
)

//go:generate mockgen -source ports.go -destination ../mocks/mockVehicle.go

type VehicleService interface {
	GetVehicle(
		where map[string]string,
		orderBy map[string]bool,
		limit int,
		offset int,
	) ([]domain.Vehicle, *errs.AppError)
}

type VehicleRepository interface {
	GetVehicle(
		whereClauses []domain.WhereClause,
		orderByClauses []domain.OrderByClause,
		pagination domain.Pagination,
	) ([]domain.Vehicle, *errs.AppError)
}
