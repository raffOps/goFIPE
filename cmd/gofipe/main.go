package main

import (
	"github.com/rjribeiro/goFIPE/cmd/gofipe/database"
	"github.com/rjribeiro/goFIPE/cmd/gofipe/routes"
	"os"
)

func main() {
	postresHost := os.Getenv("POSTGRES_HOST")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")
	database.ConnectToDB(
		postresHost,
		postgresUser,
		postgresPassword,
		postgresDB,
	)
	r := routes.SetupRouter()
	_ = r.Run(":8080")
}
