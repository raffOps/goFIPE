package useCases

import (
	"fmt"
	"github.com/rjribeiro/goFIPE/cmd/gofipe/models"
	"gorm.io/gorm"
	"time"
)

const (
	MaxLimit = 300
)

type Fipe struct {
	db *gorm.DB
}

func (v *Fipe) getVehicles(conditions string,
	args []interface{},
	orderBy string,
	offset uint,
	limit uint) ([]models.Vehicle, error) {

	if limit < 1 || limit > MaxLimit {
		return nil, &InvalidLimitError{limit: limit}
	}

	var vehicles []models.Vehicle
	err := v.db.
		Omit("ID", "CreatedAt", "UpdatedAt", "DeletedAt").
		Where(conditions, args...).
		Order(orderBy).
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&vehicles).Error
	return vehicles, err
}

func (v *Fipe) GetVehiclesByFipeCode(fipeCode string, offset uint, limit uint) ([]models.Vehicle, error) {
	return v.getVehicles("fipe_code = ?",
		[]interface{}{fipeCode},
		"year asc, month asc",
		offset,
		limit,
	)
}

func (v *Fipe) isValidReferenceMonth(reference ReferenceMonth) error {
	date := fmt.Sprintf("%d-%02d-01", reference.Year, reference.Month)
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return &InvalidReferenceMonthError{reference}
	}
	if parsedDate.After(time.Now()) {
		return &ReferenceMonthInTheFutureError{reference}
	}
	return nil
}

type ReferenceMonth struct {
	Year  uint
	Month uint
}

func (v *Fipe) GetVehiclesByReferenceMonth(reference ReferenceMonth, offset uint, limit uint) ([]models.Vehicle, error) {
	err := v.isValidReferenceMonth(reference)
	if err != nil {
		return nil, err
	}

	return v.getVehicles("year = ? and month = ?",
		[]interface{}{reference.Year, reference.Month}, "fipe_code asc",
		offset,
		limit,
	)
}

// InsertVehicle is a helper function to extract the IDs from the vehicles
func (v *Fipe) InsertVehicle(vehicles []models.Vehicle) ([]uint, error) {
	res := v.db.Create(vehicles)
	if res.Error != nil || res.RowsAffected != int64(len(vehicles)) {
		return nil, res.Error
	}
	return getVehicleIDs(vehicles), nil
}

// Extract Method applied to fetch vehicle IDs
func getVehicleIDs(vehicles []models.Vehicle) []uint {
	var ids []uint
	for _, vehicle := range vehicles {
		ids = append(ids, vehicle.ID)
	}
	return ids
}
