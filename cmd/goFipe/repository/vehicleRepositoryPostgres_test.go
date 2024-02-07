package repository

import (
	"fmt"
	"github.com/ory/dockertest/v3"
	postgresDb "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"reflect"
	"testing"
)

func init() {
	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	network, err := pool.CreateNetwork("backend")
	if err != nil {
		log.Fatalf("Could not create Network to docker: %s \n", err)
	}
	_, err = startPostgres(pool, network)
	if err != nil {
		pool.RemoveNetwork(network)
	}
}

func startPostgres(pool *dockertest.Pool, network *dockertest.Network) (*dockertest.Resource, error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "postgres",
		Repository: "postgres",
		Tag:        "13",
		Networks:   []*dockertest.Network{network},
		Env: []string{
			"POSTGRES_USER=test",
			"POSTGRES_PASSWORD=test",
			"POSTGRES_DB=test",
		},
	})

	resource.GetPort("5432/tcp")
	err = testPostgresConnection()

	if err != nil {
		resource.Close()
		return nil, err

	}

	return resource, nil
}

func testPostgresConnection() error {
	var err error
	var db *gorm.DB
	host := "localhost"
	user := "test"
	password := "test"
	database := "test"
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		user,
		password,
		database,
	)
	db, err2 := gorm.Open(
		postgresDb.Open(dsn),
		&gorm.Config{},
	)
	if err2 != nil {
		return err2
	}

	sqlDB, err3 := db.DB()
	if err != nil {
		return err3
	}

	// Ping function checks the database connectivity
	err4 := sqlDB.Ping()
	if err != nil {
		return err4
	}

	fmt.Println("Connected to the database!")
	return nil
}

func TestNewVehicleRepositoryPostgres(t *testing.T) {
	type args struct {
		conn *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *VehicleRepositoryPostgres
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVehicleRepositoryPostgres(tt.args.conn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVehicleRepositoryPostgres() = %v, want %v", got, tt.want)
			}
		})
	}
}
