package service

import (
	"fmt"
	"github.com/raffops/gofipe/cmd/goFipe/domain"
	"github.com/raffops/gofipe/cmd/goFipe/errs"
	"github.com/raffops/gofipe/cmd/goFipe/logger"
	"github.com/raffops/gofipe/cmd/goFipe/port"
)

type VehicleService struct {
	vehicleRepo port.VehicleRepository
}

func NewVehicleService(vehicleRepo port.VehicleRepository) VehicleService {
	return VehicleService{vehicleRepo: vehicleRepo}
}

func (v VehicleService) GetVehicle(
	where map[string]string,
	orderBy map[string]bool,
	offset int,
	limit int) ([]domain.Vehicle, *errs.AppError) {

	logger.Info("GetVehicle service called",

		logger.String("where", fmt.Sprint(where)),
		logger.String("orderBy", fmt.Sprint(orderBy)),
		logger.Int("offset", offset),
		logger.Int("limit", limit),
	)

	var whereClauses []domain.WhereClause
	for column, value := range where {
		whereClauses = append(whereClauses, domain.WhereClause{
			Column:   column,
			Operator: "=",
			Value:    value,
		})
	}

	var orderByClauses []domain.OrderByClause
	for column, isDesc := range orderBy {
		orderByClauses = append(
			orderByClauses,
			domain.OrderByClause{Column: column, IsDesc: isDesc},
		)
	}

	if limit == 0 {
		limit = domain.MaxLimit
	}
	pagination := domain.Pagination{
		Offset: offset,
		Limit:  limit,
	}

	return v.vehicleRepo.GetVehicle(whereClauses, orderByClauses, pagination)
}
