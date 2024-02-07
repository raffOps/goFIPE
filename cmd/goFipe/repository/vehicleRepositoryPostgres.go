package repository

import (
	"errors"
	"fmt"
	"github.com/raffops/gofipe/cmd/goFipe/domain"
	"github.com/raffops/gofipe/cmd/goFipe/errs"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

const MaxLimit = 100

type VehicleRepositoryPostgres struct {
	conn *gorm.DB
}

type Vehicle struct {
	*gorm.Model
	Year           int `gorm:"index:idx_year"`
	Month          int `gorm:"index:idx_month"`
	FipeCode       int `gorm:"index:idx_fipe_code"`
	Brand          string
	VehicleModel   string
	YearModel      string
	Authentication string
	MeanValue      float64
}

func NewVehicleRepositoryPostgres(conn *gorm.DB) *VehicleRepositoryPostgres {
	err := conn.AutoMigrate(&Vehicle{})
	if err != nil {
		panic(err)
	}
	return &VehicleRepositoryPostgres{conn: conn}
}

// GetVehicle returns a list of vehicles based on the conditions, order by and pagination
func (v VehicleRepositoryPostgres) GetVehicle(
	conditions []domain.Condition,
	orderBy []domain.OrderBy,
	pagination domain.Pagination) ([]domain.Vehicle, *errs.AppError) {
	if err := validatePagination(pagination); err != nil {
		return nil, err
	}

	whereString, whereValues, err := buildWhereQueryAndValues(conditions)
	if err != nil {
		return nil, err
	}

	orderByString, err := buildOrderByQuery(orderBy)
	if err != nil {
		return nil, err
	}

	return fetchVehiclesFromDb(v, whereString, whereValues, orderByString)
}

func validatePagination(pagination domain.Pagination) *errs.AppError {
	if pagination.Limit < 1 || pagination.Limit > MaxLimit {
		return errs.NewInvalidLimitError(fmt.Sprintf("invalid limit. The limit must be between 1 and %d", MaxLimit))
	}
	if pagination.Offset < 0 {
		return errs.NewInvalidOffsetError("invalid offset. The offset must be greater than 0")
	}
	return nil
}

func buildWhereQueryAndValues(conditions []domain.Condition) (string, []interface{}, *errs.AppError) {
	whereString, err := buildWhereQuery(conditions)
	if err != nil {
		return "", nil, err
	}
	whereValues := make([]interface{}, len(conditions))
	for i, condition := range conditions {
		whereValues[i] = condition.Value
	}
	return whereString, whereValues, nil
}

func fetchVehiclesFromDb(v VehicleRepositoryPostgres, whereColumns string, whereValues []interface{}, orderByString string) ([]domain.Vehicle, *errs.AppError) {
	var vehicles []domain.Vehicle
	err := v.conn.
		Omit("ID", "CreatedAt", "UpdatedAt", "DeletedAt").
		Where(whereColumns, whereValues).
		Order(orderByString).
		Find(&vehicles).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFoundError("Vehicles not found for that fipe code")
		}
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	return vehicles, nil
}

func buildOrderByQuery(by []domain.OrderBy) (string, *errs.AppError) {
	var orderByQuery strings.Builder
	for _, order := range by {
		if isFieldOfStruct(domain.Vehicle{}, order.Column) {
			orderByQuery.WriteString(fmt.Sprintf("%s %s", order.Column, order.Order))
		}
		return "", errs.NewValidationError(fmt.Sprintf("invalid column %s", order.Column))
	}
	return orderByQuery.String(), nil
}

func buildWhereQuery(conditions []domain.Condition) (string, *errs.AppError) {
	var whereQuery strings.Builder

	for _, condition := range conditions {
		if isFieldOfStruct(domain.Vehicle{}, condition.Column) {
			whereQuery.WriteString(fmt.Sprintf("%s = ?", condition.Column))
		}
		return "", errs.NewValidationError(fmt.Sprintf("invalid column %s", condition.Column))
	}
	return whereQuery.String(), nil
}

func isFieldOfStruct(stru interface{}, fieldName string) bool {
	v := reflect.ValueOf(stru)
	for i := 0; i < v.NumField(); i++ {
		if fieldName == v.Type().Field(i).Name {
			return true
		}
	}
	return false
}
