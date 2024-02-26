package main

import (
	"github.com/raffops/gofipe/cmd/goFipe/controller/rest"
	"github.com/raffops/gofipe/cmd/goFipe/database/postgres"
	"github.com/raffops/gofipe/cmd/goFipe/logger"
	postgresRepo "github.com/raffops/gofipe/cmd/goFipe/repository/postgres"
	"github.com/raffops/gofipe/cmd/goFipe/service"
	"github.com/raffops/gofipe/cmd/goFipe/utils"
)

func main() {
	logger.Info("Starting api...")
	errEnvVariables := utils.LoadEnvVariables()
	if errEnvVariables != nil {
		logger.Fatal("Error loading environment variables", logger.String("error", errEnvVariables.Message))
	}
	postgresConn := postgres.GetPostgresConnection()
	defer postgres.ClosePostgresConnection(postgresConn)

	vehicleRepo := postgresRepo.NewVehicleRepositoryPostgres(postgresConn)
	vehicleService := service.NewVehicleService(vehicleRepo)
	rest.Start(vehicleService)
}
