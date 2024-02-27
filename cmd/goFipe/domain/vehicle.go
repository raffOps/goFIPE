package domain

import (
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type Vehicle struct {
	Year           int    `validate:"required"`
	Month          int    `validate:"required"`
	FipeCode       string `validate:"required,validateFipeCode"`
	Brand          string
	Model          string
	YearModel      string
	Authentication string
	MeanValue      float32
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

	if IsValidYearMonth(vehicle.Year, vehicle.Month) {
		sl.ReportError(vehicle.Year, "Ano", "Year", "yearfuture", "")
		sl.ReportError(vehicle.Month, "Mes", "Month", "monthfuture", "")
	}
}

// IsValidYearMonth validates if the year/month are in the future or with invalid values
func IsValidYearMonth(year int, month int) bool {
	currentYear, currentMonth, _ := time.Now().Date()
	if year > currentYear {
		return false
	}
	if year == currentYear && month > int(currentMonth) {
		return false
	}
	if year < 1900 || month < 1 || month > 12 {
		return false
	}
	return true
}

func IsValidYear(year int) bool {
	currentYear, _, _ := time.Now().Date()
	return year >= 1900 && year <= currentYear
}

func IsValidMonth(month int) bool {
	return month >= 1 && month <= 12
}

// validateFipeCode validates if the fipe code is in the correct format
func validateFipeCode(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return IsValidFipeCode(field)
}

// IsValidFipeCode validates if the fipe code is in the correct format
func IsValidFipeCode(fipeCode string) bool {
	matched, _ := regexp.Match("^[0-9]{6}-[0-9]$", []byte(fipeCode))
	return matched
}

func GetDomainVehiclesExamples() []Vehicle {
	return []Vehicle{
		{
			Year:           2021,
			Month:          7,
			FipeCode:       "111111-1",
			Brand:          "Acura",
			Model:          "Integra GS 1.8",
			YearModel:      "1992 Gasolina",
			Authentication: "1",
			MeanValue:      700,
		},
		{
			Year:           2021,
			Month:          6,
			FipeCode:       "222222-2",
			Brand:          "Fiat",
			Model:          "147 C/ CL",
			YearModel:      "1991 Gasolina",
			Authentication: "2",
			MeanValue:      800,
		},
		{
			Year:           2021,
			Month:          7,
			FipeCode:       "222222-2",
			Brand:          "Fiat",
			Model:          "147 C/ CL",
			YearModel:      "1991 Gasolina",
			Authentication: "2",
			MeanValue:      801,
		},
		{
			Year:           2021,
			Month:          8,
			FipeCode:       "333333-3",
			Brand:          "Fiat",
			Model:          "147 C/ CL",
			YearModel:      "1991 Gasolina",
			Authentication: "2",
			MeanValue:      802,
		},
	}
}
