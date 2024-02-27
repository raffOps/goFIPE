package postgres

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/raffops/gofipe/cmd/goFipe/domain"
	"github.com/raffops/gofipe/cmd/goFipe/errs"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VehicleRepositoryPostgres struct {
	Conn *gorm.DB
}

type Vehicle struct {
	Year           int     `gorm:"index:idx_year" json:"year,omitempty"`
	Month          int     `gorm:"index:idx_month" json:"month,omitempty"`
	FipeCode       string  `gorm:"index:idx_fipe_code" json:"fipe_code,omitempty"`
	Brand          string  `json:"brand,omitempty"`
	VehicleModel   string  `json:"vehicle_model,omitempty"`
	YearModel      string  `json:"year_model,omitempty"`
	Authentication string  `json:"authentication,omitempty"`
	MeanValue      float32 `json:"mean_value,omitempty"`
}

// NewVehicleRepositoryPostgres initializes a new instance of VehicleRepositoryPostgres with the given database connection.
// It performs automatic migrations for the Vehicle model and panics if an error occurs during migration.
// Returns a pointer to the VehicleRepositoryPostgres instance.
func NewVehicleRepositoryPostgres(conn *gorm.DB) *VehicleRepositoryPostgres {
	err := conn.AutoMigrate(&Vehicle{})
	if err != nil {
		panic(err)
	}
	return &VehicleRepositoryPostgres{Conn: conn}
}

// GetVehicle retrieves vehicles from the database based on the given conditions, order by specifications and pagination settings.
// It validates the pagination parameters and returns an error if they are invalid.
// It then fetches the vehicles from the database using the fetchVehiclesFromDb function, passing the specified conditions, order by specifications, and pagination settings.
// If an error occurs during the fetch operation, it returns the error.
// Otherwise, it converts the fetched vehicles to domain.Vehicle objects using the ToDomainVehicles function and returns them along with a nil error.
func (v VehicleRepositoryPostgres) GetVehicle(
	whereClauses []domain.WhereClause,
	orderByClauses []domain.OrderByClause,
	pagination domain.Pagination) ([]domain.Vehicle, *errs.AppError) {
	if err := validatePagination(pagination); err != nil {
		return nil, err
	}

	vehicles, err := fetchVehiclesFromDb(v, whereClauses, orderByClauses, pagination)
	if err != nil {
		return nil, err
	}

	return ToDomainVehicles(vehicles), nil
}

func fetchVehiclesFromDb(v VehicleRepositoryPostgres,
	whereClauses []domain.WhereClause,
	orderByClauses []domain.OrderByClause,
	pagination domain.Pagination) ([]Vehicle, *errs.AppError) {

	var vehicles []Vehicle
	fetch := v.Conn.Omit("ID", "CreatedAt", "UpdatedAt", "DeletedAt")

	for _, whereClause := range whereClauses {
		if isValidJsonField(Vehicle{}, whereClause.Column) {
			query := fmt.Sprintf("%s %s ?", whereClause.Column, whereClause.Operator)
			fetch = fetch.Where(query, whereClause.Value)
		}
	}

	for _, orderByClause := range orderByClauses {
		if isValidJsonField(Vehicle{}, orderByClause.Column) {
			fetch = fetch.Order(
				clause.OrderByColumn{
					Column: clause.Column{Name: orderByClause.Column},
					Desc:   orderByClause.IsDesc,
				})
		}
	}

	fetch = fetch.Offset(pagination.Offset).Limit(pagination.Limit)

	result := fetch.Find(&vehicles)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFoundError("Vehicles not found")
		}
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	if len(vehicles) == 0 {
		return nil, errs.NewNotFoundError("Vehicles not found")
	}

	return vehicles, nil
}

// ToDomainVehicles converts a slice of Vehicle objects to a slice of domain.Vehicle objects.
func ToDomainVehicles(vehicles []Vehicle) []domain.Vehicle {
	var domainVehicles []domain.Vehicle

	for _, vehicle := range vehicles {
		domainVehicles = append(domainVehicles,
			domain.Vehicle{
				Year:           vehicle.Year,
				Month:          vehicle.Month,
				FipeCode:       vehicle.FipeCode,
				Brand:          vehicle.Brand,
				Model:          vehicle.VehicleModel,
				YearModel:      vehicle.YearModel,
				Authentication: vehicle.Authentication,
				MeanValue:      vehicle.MeanValue,
			},
		)
	}
	return domainVehicles
}

// validatePagination validates the given pagination parameters.
// It checks if the limit is within the valid range (1 to MaxLimit) and if the offset is greater than or equal to 0.
// If any validation error occurs, it returns an AppError with the corresponding error message.
// Otherwise, it returns nil indicating that the pagination is valid.
// Example usage: err := validatePagination(pagination)
func validatePagination(pagination domain.Pagination) *errs.AppError {
	if pagination.Limit < 1 || pagination.Limit > domain.MaxLimit {
		return errs.NewUnprocessableEntityError(
			fmt.Sprintf("invalid limit. The limit must be between 1 and %d",
				domain.MaxLimit,
			),
		)

	}
	if pagination.Offset < 0 {
		return errs.NewUnprocessableEntityError("invalid offset. The offset must be greater than 0")
	}
	if pagination.Offset > pagination.Limit {
		return errs.NewUnprocessableEntityError("Offset must be smaller than Limit")

	}
	return nil
}

// isValidJsonField returns true if the given column is a valid JSON field in the input struct, otherwise false.
func isValidJsonField(input interface{}, column string) bool {
	typeOfInput := reflect.TypeOf(input)

	for i := 0; i < typeOfInput.NumField(); i++ {
		// GetByFipeCode the field, returns https://golang.org/pkg/reflect/#StructField
		field := typeOfInput.Field(i)
		jsonTag, ok := field.Tag.Lookup("json")
		jsonTagList := strings.Split(jsonTag, ",")
		if ok && slices.Contains(jsonTagList, column) {
			return true
		}
	}
	return false
}
