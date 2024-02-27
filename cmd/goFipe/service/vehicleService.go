package service

import (
	"fmt"
	"github.com/raffops/gofipe/cmd/goFipe/domain"
	"github.com/raffops/gofipe/cmd/goFipe/errs"
	"github.com/raffops/gofipe/cmd/goFipe/logger"
	"github.com/raffops/gofipe/cmd/goFipe/port"
	"slices"
	"strconv"
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

	errValidate := validate(where, orderBy, offset, limit)
	if errValidate != nil {
		return nil, errValidate
	}

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

	pagination := domain.Pagination{
		Offset: offset,
		Limit:  limit,
	}

	return v.vehicleRepo.GetVehicle(whereClauses, orderByClauses, pagination)
}

func validate(where map[string]string, orderBy map[string]bool, offset int, limit int) *errs.AppError {
	if errValidateWhere := validateWhere(where); errValidateWhere != nil {
		return errValidateWhere
	}

	if errValidateOrder := validateOrderBy(orderBy); errValidateOrder != nil {
		return errValidateOrder
	}

	if errValidatePagination := validatePagination(offset, limit); errValidatePagination != nil {
		return errValidatePagination
	}
	return nil
}

func validateWhere(where map[string]string) *errs.AppError {
	if len(where) == 0 {
		return errs.NewBadRequestError("Where is required")
	}
	for column, value := range where {
		switch column {
		case "fipe_code":
			if !domain.IsValidFipeCode(value) {
				return errs.NewValidationError("Invalid fipe code")
			}
		case "year":
			value, err := strconv.Atoi(value)
			if err != nil {
				return errs.NewValidationError("Invalid year")
			}
			if !domain.IsValidYear(value) {
				return errs.NewValidationError("Invalid year")
			}
		case "month":
			value, err := strconv.Atoi(value)
			if err != nil {
				return errs.NewValidationError("Invalid month")
			}
			if !domain.IsValidMonth(value) {
				return errs.NewValidationError("Invalid month")
			}
		case "mean_value":
			_, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return errs.NewValidationError("Invalid mean value")
			}
		default:
			return errs.NewValidationError("Invalid Column")
		}
	}
	return nil
}

func validateOrderBy(orderBy map[string]bool) *errs.AppError {
	validColumns := []string{"fipe_code", "year", "month", "mean_value"}
	if len(orderBy) == 0 {
		return errs.NewBadRequestError("OrderBy is required")
	}
	for column := range orderBy {
		if !slices.Contains(validColumns, column) {
			return errs.NewValidationError(fmt.Sprintf("Invalid column: %s", column))
		}
	}
	return nil
}

func validatePagination(offset int, limit int) *errs.AppError {
	if offset < 0 {
		return errs.NewValidationError("Offset must be greater than 0")
	}
	if offset > limit {
		return errs.NewValidationError("Offset must be smaller than Limit")
	}
	if limit > domain.MaxLimit {
		return errs.NewValidationError(
			fmt.Sprintf("Limit must be smaller or equal than %d",
				domain.MaxLimit,
			),
		)
	}

	return nil
}
