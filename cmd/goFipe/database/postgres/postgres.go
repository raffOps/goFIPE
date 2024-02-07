package postgres

import (
	"fmt"
	postgresDb "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

// sanityCheck checks if the environment variables are set.
func sanityCheck() {
	if _, ok := os.LookupEnv("POSTGRES_HOST"); !ok {
		panic("POSTGRES_HOST environment variable is not set")
	}

	if _, ok := os.LookupEnv("POSTGRES_USER"); !ok {
		panic("POSTGRES_USER environment variable is not set")
	}

	if _, ok := os.LookupEnv("POSTGRES_PASSWORD"); !ok {
		panic("POSTGRES_PASSWORD environment variable is not set")
	}

	if _, ok := os.LookupEnv("POSTGRES_DB"); !ok {
		panic("POSTGRES_DB environment variable is not set")
	}
}

// GetPostgresConnection returns a PostgreSQL database connection using the environment variables:
// POSTGRES_HOST, POSTGRES_USER, POSTGRES_PASSWORD and POSTGRES_DB.
func GetPostgresConnection() *gorm.DB {
	sanityCheck()
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
	DB, err := gorm.Open(
		postgresDb.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		panic("Unable to connect to database")
	}
	return DB
}

func ClosePostgresConnection(conn *gorm.DB) {
	sqlDB, _ := conn.DB()
	_ = sqlDB.Close()
}
