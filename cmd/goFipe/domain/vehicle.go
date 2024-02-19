package domain

import (
	"github.com/go-playground/validator/v10"
	"github.com/raffops/gofipe/cmd/goFipe/errs"
	"regexp"
	"time"
)

//go:generate mockgen -source vehicle.go -destination ../mocks/domain/mockVehicle.go

type Vehicle struct {
	Year           int `validate:"required"`
	Month          int `validate:"required"`
	FipeCode       int `validate:"required,validateFipeCode"`
	Brand          string
	Model          string
	YearModel      string
	Authentication string
	MeanValue      float64
}

type VehicleService interface {
	GetVehicleByFipeCode(fipeCode int, orderBy []OrderBy, pagination Pagination) ([]Vehicle, *errs.AppError)
	GetVehicleByReferenceYearMonth(year int, month int, orderBy []OrderBy, pagination Pagination) ([]Vehicle, *errs.AppError)
	InsertVehicles([]Vehicle, *errs.AppError)
}

type VehicleRepository interface {
	GetVehicle(conditions []Condition, orderBy []OrderBy, pagination Pagination) ([]Vehicle, *errs.AppError)
	InsertVehicles([]Vehicle) *errs.AppError
}

// Validate validates the vehicle struct
func (v *Vehicle) Validate() error {
	validate := validator.New()
	_ = validate.RegisterValidation("validateFipeCode", validateFipeCode)
	validate.RegisterStructValidation(validateYearMonth, Vehicle{})
	return validate.Struct(v)
}

// validateYearMonth validates if the year/month are in the future
func validateYearMonth(sl validator.StructLevel) {
	vehicle := sl.Current().Interface().(Vehicle)

	currentYear, currentMonth, _ := time.Now().Date()
	if vehicle.Year > currentYear || (vehicle.Year == currentYear && vehicle.Month > int(currentMonth)) {
		sl.ReportError(vehicle.Year, "Ano", "Year", "yearfuture", "")
		sl.ReportError(vehicle.Month, "Mes", "Month", "monthfuture", "")
	}
}

// validateFipeCode validates if the fipe code is in the correct format
func validateFipeCode(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	matched, _ := regexp.Match("^[0-9]{6}-[0-9]$", []byte(field))
	return matched
}

func GetDomainVehiclesExamples() []Vehicle {
	return []Vehicle{
		{
			Year:           2021,
			Month:          7,
			FipeCode:       1,
			Brand:          "Acura",
			Model:          "Integra GS 1.8",
			YearModel:      "1992 Gasolina",
			Authentication: "1",
			MeanValue:      700,
		},
		{
			Year:           2021,
			Month:          6,
			FipeCode:       2,
			Brand:          "Fiat",
			Model:          "147 C/ CL",
			YearModel:      "1991 Gasolina",
			Authentication: "2",
			MeanValue:      800,
		},
		{
			Year:           2021,
			Month:          7,
			FipeCode:       2,
			Brand:          "Fiat",
			Model:          "147 C/ CL",
			YearModel:      "1991 Gasolina",
			Authentication: "2",
			MeanValue:      801,
		},
		{
			Year:           2021,
			Month:          8,
			FipeCode:       2,
			Brand:          "Fiat",
			Model:          "147 C/ CL",
			YearModel:      "1991 Gasolina",
			Authentication: "2",
			MeanValue:      802,
		},
	}
}
