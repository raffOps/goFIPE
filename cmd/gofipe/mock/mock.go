package mock

import (
	"github.com/rjribeiro/goFIPE/cmd/gofipe/database"
	"github.com/rjribeiro/goFIPE/cmd/gofipe/models"
)

// GetVehiclesOnDB return mocked vehicles on database
func GetVehiclesOnDB() []models.Vehicle {
	vehicles := []models.Vehicle{
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
			Value:          80876.0,
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
			Month:          10,
			Value:          17613,
			Brand:          "Walk",
			ModelCar:       "Buggy  Walk Sport 1.6 8V 58cv",
			YearModel:      "2005",
			Fuel:           "Gasolina",
			FipeCode:       "061001-1",
			VehicleType:    1,
			FuelAcronym:    "G",
			ExtractionDate: "2023-12-17 17:23:17",
			FlagTest:       true,
		},
		{
			Year:           2023,
			Month:          11,
			Value:          17513,
			Brand:          "Walk",
			ModelCar:       "Buggy  Walk Sport 1.6 8V 58cv",
			YearModel:      "2005",
			Fuel:           "Gasolina",
			FipeCode:       "061001-1",
			VehicleType:    1,
			FuelAcronym:    "G",
			ExtractionDate: "2023-12-17 17:23:17",
			FlagTest:       true,
		},
	}
	for _, vehicle := range vehicles {
		database.DB.Create(&vehicle)
	}
	return vehicles
}
