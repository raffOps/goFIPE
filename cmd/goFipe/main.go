package main

import (
	"github.com/raffops/gofipe/cmd/goFipe/controller/rest"
	"github.com/raffops/gofipe/cmd/goFipe/database/postgres"
	postgresRepo "github.com/raffops/gofipe/cmd/goFipe/repository/postgres"
	"github.com/raffops/gofipe/cmd/goFipe/service"
)

func main() {
	postgresConn := postgres.GetPostgresConnection()
	defer postgres.ClosePostgresConnection(postgresConn)

	vehicleRepo := postgresRepo.NewVehicleRepositoryPostgres(postgresConn)
	vehicleService := service.NewVehicleService(vehicleRepo)
	rest.Start(vehicleService)
}
