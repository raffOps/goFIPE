package database

import (
	"fmt"
	"github.com/rjribeiro/goFIPE/cmd/gofipe/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	DB  *gorm.DB
	err error
)

func ConnectToDB() {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	database := os.Getenv("POSTGRES_DB")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=UTC",
		host,
		user,
		password,
		database,
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("Unable to connect to database")
	}
	err = DB.AutoMigrate(&models.Test{}, &models.Vehicle{})
	if err != nil {
		log.Panic(err)
	}
}
