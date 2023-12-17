package main

import (
	"github.com/rjribeiro/goFIPE/cmd/gofipe/database"
	"github.com/rjribeiro/goFIPE/cmd/gofipe/routes"
)

func main() {

	database.ConnectToDB()
	r := routes.SetupRouter()
	_ = r.Run(":8080")
}
