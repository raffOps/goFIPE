package useCases

import (
	"github.com/gin-gonic/gin"
	"github.com/rjribeiro/goFIPE/cmd/gofipe/database"
	"github.com/rjribeiro/goFIPE/cmd/gofipe/mock"
	"github.com/rjribeiro/goFIPE/cmd/gofipe/models"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func init() {
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_USER", "root")
	os.Setenv("POSTGRES_PASSWORD", "root")
	os.Setenv("POSTGRES_DB", "gofipe_test")
	database.ConnectToDB()
	gin.SetMode(gin.TestMode)
}

type TestCaseGetVehiclesByFipeCode struct {
	name     string
	fipeCode string
	offset   uint
	limit    uint
	err      error
	wanted   []models.Vehicle
}

// GetTestCasesGetVehiclesByFipeCode generates a list of test cases for the GetVehiclesByFipeCode function.
// Each test case includes the fipeCode fipeCode, offset, limit, expected error, and expected list of vehicles.
// The test cases cover different scenarios such as valid fipeCode, different offset and limit values, and empty result.
func GetTestCasesGetVehiclesByFipeCode(vehiclesOnDb []models.Vehicle) []TestCaseGetVehiclesByFipeCode {
	testCases := []TestCaseGetVehiclesByFipeCode{
		{ // Test Case 1
			name:     "Valid fipeCode, offset 0, limit 1",
			fipeCode: vehiclesOnDb[0].FipeCode,
			offset:   0,
			limit:    1,
			err:      nil,
			wanted:   []models.Vehicle{vehiclesOnDb[0]},
		},
		{
			name:     "Valid fipeCode, offset 1, limit 1",
			fipeCode: vehiclesOnDb[1].FipeCode,
			offset:   1,
			limit:    1,
			err:      nil,
			wanted:   []models.Vehicle{vehiclesOnDb[1]},
		},
		{
			name:     "Valid fipeCode, offset 0, limit 2",
			fipeCode: vehiclesOnDb[0].FipeCode,
			offset:   0,
			limit:    2,
			err:      nil,
			wanted:   []models.Vehicle{vehiclesOnDb[0], vehiclesOnDb[1]},
		},
		{
			name:     "Valid fipeCode, offset 1, limit 2, empty result",
			fipeCode: "653433-2",
			offset:   4,
			limit:    2,
			err:      nil,
			wanted:   []models.Vehicle{},
		},
	}

	return testCases
}

func TestFipe_GetVehiclesByFipeCode(t *testing.T) {
	vehiclesOnDb := mock.GetVehiclesOnDB()
	t.Cleanup(func() { database.TruncateTable("vehicles") })

	testCases := GetTestCasesGetVehiclesByFipeCode(vehiclesOnDb)

	fipeTable := Fipe{db: database.DB}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			vehicles, err := fipeTable.GetVehiclesByFipeCode(tt.fipeCode, tt.offset, tt.limit)
			assert.Equal(t, tt.err, err)
			if tt.err == nil {
				assert.Equal(t, tt.wanted, vehicles)
			}
		},
		)
	}
}

type TestCaseGetVehiclesByReferenceMonth struct {
	name      string
	reference ReferenceMonth
	offset    uint
	limit     uint
	err       error
	wanted    []models.Vehicle
}

// GetTestCasesGetVehiclesByReferenceMonth generates a list of test cases for the GetVehiclesByReferenceMonth function.
// Each test case includes the fipeCode year, month, offset, limit, expected error, and expected list of vehicles.
// The test cases cover different scenarios such as valid year/month, different offset and limit values, and empty result.
func GetTestCasesGetVehiclesByReferenceMonth(vehiclesOnDb []models.Vehicle) []TestCaseGetVehiclesByReferenceMonth {
	testCases := []TestCaseGetVehiclesByReferenceMonth{
		{
			name:      "2023-11, offset 0, limit 1",
			reference: ReferenceMonth{2023, 11},
			offset:    0,
			limit:     1,
			err:       nil,
			wanted:    []models.Vehicle{vehiclesOnDb[1]},
		},
		{
			name:      "2023-11, offset 1, limit 1",
			reference: ReferenceMonth{2023, 11},
			offset:    1,
			limit:     1,
			err:       nil,
			wanted:    []models.Vehicle{vehiclesOnDb[3]},
		},
		{
			name:      "2023-10, offset 0, limit 2",
			reference: ReferenceMonth{2023, 10},
			offset:    0,
			limit:     2,
			err:       nil,
			wanted:    []models.Vehicle{vehiclesOnDb[0], vehiclesOnDb[2]},
		},
		{
			name:      "2023-01 offset 1, limit 2, empty result",
			reference: ReferenceMonth{2023, 1},
			offset:    4,
			limit:     2,
			err:       nil,
			wanted:    []models.Vehicle{},
		},
		{
			name:      "Invalid reference month",
			reference: ReferenceMonth{2023, 13},
			offset:    0,
			limit:     1,
			err:       &InvalidReferenceMonthError{ReferenceMonth{2023, 13}},
		},
		{
			name:      "Reference month in the future",
			reference: ReferenceMonth{2080, 12},
			offset:    0,
			limit:     1,
			err:       &ReferenceMonthInTheFutureError{ReferenceMonth{2080, 12}},
		},
		{
			name:      "limit less than 1",
			reference: ReferenceMonth{2023, 10},
			offset:    0,
			limit:     0,
			err:       &InvalidLimitError{limit: 0},
		},
		{
			name:      "limit greater than MaxLimit",
			reference: ReferenceMonth{2023, 10},
			offset:    0,
			limit:     MaxLimit + 1,
			err:       &InvalidLimitError{limit: MaxLimit + 1},
		},
	}
	return testCases
}

func TestFipe_GetVehiclesByReferenceMonth(t *testing.T) {
	vehiclesOnDb := mock.GetVehiclesOnDB()
	t.Cleanup(func() { database.TruncateTable("vehicles") })

	testCases := GetTestCasesGetVehiclesByReferenceMonth(vehiclesOnDb)

	fipeTable := Fipe{db: database.DB}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			vehicles, err := fipeTable.GetVehiclesByReferenceMonth(tt.reference, tt.offset, tt.limit)
			assert.Equal(t, tt.err, err)
			if tt.err == nil {
				assert.Equal(t, tt.wanted, vehicles)
			}
		},
		)
	}
}

type TestCaseInsertVehicle struct {
	name  string
	input []models.Vehicle
}

func GetTestCasesInsertVehicle() []TestCaseInsertVehicle {
	testCases := []TestCaseInsertVehicle{
		{ // Test Case 1
			name: "Insert one vehicle",
			input: []models.Vehicle{
				{
					Year:           2023,
					Month:          10,
					Value:          80877.0,
					Brand:          "Renault",
					ModelCar:       "STEPWAY Intense Flex 1.6 16V  Aut.",
					YearModel:      "2022",
					Fuel:           "Gasolina",
					FipeCode:       "025282-4",
					VehicleType:    1,
					FuelAcronym:    "G",
					ExtractionDate: "2023-12-17 17:23:17",
					FlagTest:       true,
				},
			},
		},
		{
			name: "Insert two vehicles",
			input: []models.Vehicle{
				{
					Year:           2023,
					Month:          10,
					Value:          80877.0,
					Brand:          "Renault",
					ModelCar:       "STEPWAY Intense Flex 1.6 16V  Aut.",
					YearModel:      "2022",
					Fuel:           "Gasolina",
					FipeCode:       "025282-4",
					VehicleType:    1,
					FuelAcronym:    "G",
					ExtractionDate: "2023-12-17 17:23:17",
					FlagTest:       true,
				},
				{
					Year:           2023,
					Month:          11,
					Value:          80877.0,
					Brand:          "Renault",
					ModelCar:       "STEPWAY Intense Flex 1.6 16V  Aut.",
					YearModel:      "2022",
					Fuel:           "Gasolina",
					FipeCode:       "025282-4",
					VehicleType:    1,
					FuelAcronym:    "G",
					ExtractionDate: "2023-12-17 17:23:17",
					FlagTest:       true,
				},
			},
		},
	}

	return testCases
}

func TestFipe_InsertVehicle(t *testing.T) {

	t.Cleanup(func() { database.TruncateTable("vehicles") })

	testCases := GetTestCasesInsertVehicle()

	fipeTable := Fipe{db: database.DB}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ids, _ := fipeTable.InsertVehicle(tt.input)
			assert.Equal(t, len(tt.input), len(ids))
		},
		)
	}
}
