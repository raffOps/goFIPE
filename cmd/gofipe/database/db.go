package database

import (
	"fmt"
	"github.com/rjribeiro/goFIPE/cmd/gofipe/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var (
	DB  *gorm.DB
	err error
)

func ConnectToDB(host string, user string, password string, database string) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=UTC",
		host,
		user,
		password,
		database,
	)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("Nao foi possivel conectar o banco de dados")
	}
	err = DB.AutoMigrate(&models.Teste{}, &models.Veiculo{})
	if err != nil {
		log.Panic(err)
	}
}
