package postgres

import (
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	postgres2 "github.com/raffops/gofipe/cmd/goFipe/database/postgres"
	"github.com/raffops/gofipe/cmd/goFipe/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/raffops/gofipe/cmd/goFipe/domain"
	"github.com/raffops/gofipe/cmd/goFipe/errs"
	"github.com/stretchr/testify/assert"
)

func Test_isValidField(t *testing.T) {
	type args struct {
		input  interface{}
		column string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "struct with valid field #1",
			args: args{
				input:  Vehicle{},
				column: "year_model",
			},
			want: true,
		},
		{
			name: "struct with valid field #2",
			args: args{
				input:  Vehicle{},
				column: "fipe_code",
			},
			want: true,
		},
		{
			name: "struct without valid field",
			args: args{
				input:  Vehicle{},
				column: "fipecode",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, isValidJsonField(tt.args.input, tt.args.column), "isValidJsonField(%v, %v)", tt.args.input, tt.args.column)
		})
	}
}

func Test_validatePagination(t *testing.T) {
	type args struct {
		pagination domain.Pagination
	}
	tests := []struct {
		name string
		args args
		want *errs.AppError
	}{
		{
			name: "Valid Pagination",
			args: args{
				pagination: domain.Pagination{Offset: 0, Limit: 10},
			},
			want: nil,
		},
		{
			name: "Offset greater than limit",
			args: args{
				pagination: domain.Pagination{Offset: 11, Limit: 10},
			},
			want: errs.NewUnprocessableEntityError("Offset must be smaller than Limit"),
		},

		{
			name: "Negative Offset",
			args: args{
				pagination: domain.Pagination{Offset: -1, Limit: 10},
			},
			want: errs.NewUnprocessableEntityError("invalid offset. The offset must be greater than 0"), // Assuming this error type
		},
		{
			name: "Negative Limit",
			args: args{
				pagination: domain.Pagination{Offset: 0, Limit: -1},
			},
			want: errs.NewUnprocessableEntityError(
				fmt.Sprintf("invalid limit. The limit must be between 1 and %d",
					domain.MaxLimit),
			),
		},
		{
			name: "Zero Limit",
			args: args{
				pagination: domain.Pagination{Offset: 0, Limit: 0},
			},
			want: errs.NewUnprocessableEntityError(
				fmt.Sprintf("invalid limit. The limit must be between 1 and %d",
					domain.MaxLimit,
				),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, validatePagination(tt.args.pagination), "validatePagination(%v)", tt.args.pagination)
		})
	}
}

func TestMain(m *testing.M) {
	var errPool error
	pool, errPool := dockertest.NewPool("")
	if errPool != nil {
		logger.Error(fmt.Sprintf("Could not connect to docker: %s", errPool))
	}

	network, errNetwork := pool.CreateNetwork("backend")
	if errNetwork != nil {
		logger.Error(fmt.Sprintf("Could not create Network to docker: %s \n", errPool))
	}

	resource, errResource := startPostgres(pool, network)

	if errResource != nil {
		logger.Error(fmt.Sprintf("Could not create Postgres: %s \n", errResource))
	}

	exitCode := m.Run()
	err := TearDown(pool, network, resource)
	if err != nil {
		log.Fatalf("Could not purge resource: %v", err)
	}

	os.Exit(exitCode)
}

// TearDown purges the resources and removes the network.
func TearDown(pool *dockertest.Pool, network *dockertest.Network, resource *dockertest.Resource) error {
	if err := pool.Purge(resource); err != nil {
		return fmt.Errorf("could not purge resource: %v", err)
	}

	if err := pool.RemoveNetwork(network); err != nil {
		return fmt.Errorf("could not remove network: %v", err)
	}

	return nil
}

func startPostgres(pool *dockertest.Pool, network *dockertest.Network) (*dockertest.Resource, error) {
	resource, _ := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "postgres",
		Repository: "postgres",
		Tag:        "13",
		Networks:   []*dockertest.Network{network},
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", os.Getenv("POSTGRES_USER")),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", os.Getenv("POSTGRES_PASSWORD")),
			fmt.Sprintf("POSTGRES_DB=%s", os.Getenv("POSTGRES_DB")),
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {{HostIP: os.Getenv("POSTGRES_HOST"), HostPort: "5432"}},
		},
	})

	err := testPostgresConnection(pool)

	if err != nil {
		return nil, err
	}

	return resource, nil
}

func testPostgresConnection(pool *dockertest.Pool) error {
	var err error
	var db *gorm.DB
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

	err = pool.Retry(func() error {
		_, err2 := gorm.Open(
			postgres.Open(dsn),
			&gorm.Config{},
		)
		return err2
	})
	if err != nil {
		return err
	}

	db, _ = gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)

	sqlDB, err3 := db.DB()
	if err != nil {
		return err3
	}

	// Ping function checks the database connectivity
	err4 := sqlDB.Ping()
	if err4 != nil {
		return err4
	}

	logger.Info("Connected to the database!")
	return nil
}

func TestNewVehicleRepositoryPostgres(t *testing.T) {
	type args struct {
		conn *gorm.DB
	}
	conn := postgres2.GetPostgresConnection()
	t.Cleanup(func() { postgres2.ClosePostgresConnection(conn) })
	tests := []struct {
		name string
		args args
		want *VehicleRepositoryPostgres
	}{
		{
			name: "Single test",
			args: args{conn: conn},
			want: &VehicleRepositoryPostgres{Conn: conn},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewVehicleRepositoryPostgres(tt.args.conn)
			assert.Equal(t, got, tt.want)
		})
	}
}

func insertVehiclesOnDB(conn *gorm.DB) ([]Vehicle, []domain.Vehicle) {
	domainVehicles := domain.GetDomainVehiclesExamples()
	var vehicles []Vehicle
	for _, domainVehicle := range domainVehicles {
		fipeCode := strconv.Itoa(domainVehicle.FipeCode)
		vehicle := Vehicle{
			Year:           domainVehicle.Year,
			Month:          domainVehicle.Month,
			FipeCode:       fipeCode,
			Brand:          domainVehicle.Brand,
			VehicleModel:   domainVehicle.Model,
			YearModel:      domainVehicle.YearModel,
			Authentication: domainVehicle.Authentication,
			MeanValue:      domainVehicle.MeanValue,
		}
		conn.Create(&vehicle)
		vehicles = append(vehicles, vehicle)
	}
	return vehicles, domainVehicles
}

func TestVehicleRepositoryPostgres_GetVehicle(t *testing.T) {
	type fields struct {
		conn *gorm.DB
	}
	type args struct {
		conditions []domain.WhereClause
		orderBy    []domain.OrderByClause
		pagination domain.Pagination
	}

	conn := postgres2.GetPostgresConnection()
	t.Cleanup(func() { postgres2.ClosePostgresConnection(conn) })
	conn.AutoMigrate(&Vehicle{})
	_, domainVehiclesOnDb := insertVehiclesOnDB(conn)

	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []domain.Vehicle
		wantError *errs.AppError
	}{
		{
			name:   "One Vehicle",
			fields: fields{conn: conn},
			args: args{
				conditions: []domain.WhereClause{},
				orderBy:    []domain.OrderByClause{},
				pagination: domain.Pagination{Offset: 0, Limit: 1},
			},
			want:      []domain.Vehicle{domainVehiclesOnDb[0]},
			wantError: nil,
		},
		{
			name:   "Order by fipe code",
			fields: fields{conn: conn},
			args: args{
				conditions: []domain.WhereClause{},
				orderBy: []domain.OrderByClause{
					{
						Column: "fipe_code",
						IsDesc: false,
					},
				},
				pagination: domain.Pagination{Offset: 0, Limit: 10},
			},
			want: []domain.Vehicle{
				domainVehiclesOnDb[0],
				domainVehiclesOnDb[1],
				domainVehiclesOnDb[2],
				domainVehiclesOnDb[3],
			},
			wantError: nil,
		},
		{
			name:   "Order desc by mean value",
			fields: fields{conn: conn},
			args: args{
				conditions: []domain.WhereClause{},
				orderBy: []domain.OrderByClause{
					{
						Column: "mean_value",
						IsDesc: true,
					},
				},
				pagination: domain.Pagination{Offset: 0, Limit: 10},
			},
			want: []domain.Vehicle{
				domainVehiclesOnDb[3],
				domainVehiclesOnDb[2],
				domainVehiclesOnDb[1],
				domainVehiclesOnDb[0],
			},
			wantError: nil,
		},
		{
			name:   "Order asc by mean value offset 1 limit 2",
			fields: fields{conn: conn},
			args: args{
				conditions: []domain.WhereClause{},
				orderBy: []domain.OrderByClause{
					{
						Column: "mean_value",
						IsDesc: false,
					},
				},
				pagination: domain.Pagination{Offset: 1, Limit: 2},
			},
			want: []domain.Vehicle{
				domainVehiclesOnDb[1],
				domainVehiclesOnDb[2],
			},
			wantError: nil,
		},
		{
			name:   "fipe_code equal to 1",
			fields: fields{conn: conn},
			args: args{
				conditions: []domain.WhereClause{
					{
						Column:   "fipe_code",
						Operator: "=",
						Value:    "1",
					},
				},
				orderBy: []domain.OrderByClause{
					{
						Column: "mean_value",
						IsDesc: false,
					},
				},
				pagination: domain.Pagination{Offset: 0, Limit: 10},
			},
			want: []domain.Vehicle{
				domainVehiclesOnDb[0],
			},
			wantError: nil,
		},
		{
			name:   "NewNotFoundError",
			fields: fields{conn: conn},
			args: args{
				conditions: []domain.WhereClause{
					{
						Column:   "fipe_code",
						Operator: "=",
						Value:    "76576",
					},
				},
				orderBy:    []domain.OrderByClause{},
				pagination: domain.Pagination{Offset: 0, Limit: 10},
			},
			want:      nil,
			wantError: errs.NewNotFoundError("Vehicles not found"),
		},
		{
			name:   "year equal to 2021 and month equal to 8",
			fields: fields{conn: conn},
			args: args{
				conditions: []domain.WhereClause{
					{
						Column:   "year",
						Operator: "=",
						Value:    2021,
					},
					{
						Column:   "month",
						Operator: "=",
						Value:    8,
					},
				},
				orderBy:    []domain.OrderByClause{},
				pagination: domain.Pagination{Offset: 0, Limit: 10},
			},
			want: []domain.Vehicle{
				domainVehiclesOnDb[3],
			},
			wantError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := VehicleRepositoryPostgres{
				Conn: tt.fields.conn,
			}
			got, gotError := v.GetVehicle(tt.args.conditions, tt.args.orderBy, tt.args.pagination)
			assert.Equalf(t, tt.want, got, "GetVehicle(%v, %v, %v)", tt.args.conditions, tt.args.orderBy, tt.args.pagination)
			assert.Equalf(t, tt.wantError, gotError, "GetVehicle(%v, %v, %v)", tt.args.conditions, tt.args.orderBy, tt.args.pagination)
		})
	}
}
