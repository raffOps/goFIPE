package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"regexp"
	"time"
)

// Vehicle is the struct that represents the vehicle table in the database,
// created according to fipe table. See: https://veiculos.fipe.org.br/
type Vehicle struct {
	gorm.Model     `json:"-"`
	Year           uint16  `json:"ano_referencia" validate:"gte=1900,lte=2100"`
	Month          uint8   `json:"mes_referencia" validate:"gte=1,lte=12"`
	Value          float32 `json:"valor" validate:"gt=0"`
	Brand          string  `json:"marca" gorm:"index"`
	ModelCar       string  `json:"modelo" gorm:"index"`
	YearModel      string  `json:"ano_modelo" gorm:"index"`
	Fuel           string  `json:"combustivel"`
	FipeCode       string  `json:"codigo_fipe" validate:"validateFipeCode" gorm:"index"`
	VehicleType    int     `json:"tipo_veiculo"`
	FuelAcronym    string  `json:"sigla_combustivel"`
	ExtractionDate string  `json:"data_consulta" validate:"validateExtractionDate"`
	FlagTest       bool    `json:"-" default:"false" gorm:"default:false"`
}

// Validate validates the vehicle struct
func (v *Vehicle) Validate() error {
	validate := validator.New()
	_ = validate.RegisterValidation("validateExtractionDate", validateExtractionDate)
	_ = validate.RegisterValidation("validateFipeCode", validateFipeCode)
	validate.RegisterStructValidation(validateYearMonth, Vehicle{})
	return validate.Struct(v)
}

// validateYearMonth validates if the year/month are in the future
func validateYearMonth(sl validator.StructLevel) {
	vehicle := sl.Current().Interface().(Vehicle)

	currentYear, currentMonth, _ := time.Now().Date()
	if int(vehicle.Year) > currentYear || (int(vehicle.Year) == currentYear && int(vehicle.Month) > int(currentMonth)) {
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

// validateExtractionDate validates if the date is in the correct format
func validateExtractionDate(fl validator.FieldLevel) bool {
	const layout = "2006-01-02 15:04:05"
	s := fl.Field().String()
	_, err := time.Parse(layout, s)
	return err == nil
}
