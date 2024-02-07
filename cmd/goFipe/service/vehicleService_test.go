package service

import (
	"github.com/golang/mock/gomock"
	"github.com/raffops/gofipe/cmd/goFipe/domain"
	"github.com/raffops/gofipe/cmd/goFipe/errs"
	mock_domain "github.com/raffops/gofipe/cmd/goFipe/mocks/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getMockVehicleRepository(t *testing.T) (*mock_domain.MockVehicleRepository, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockVehicleRepository := mock_domain.NewMockVehicleRepository(ctrl)
	return mockVehicleRepository, ctrl
}

func TestNewVehicleService(t *testing.T) {
	mockVehicleRepository, ctrl := getMockVehicleRepository(t)
	t.Cleanup(ctrl.Finish)

	type args struct {
		vehicleRepo domain.VehicleRepository
	}

	tests := []struct {
		name string
		args args
		want VehicleService
	}{
		{
			name: "TestNewVehicleService",
			args: args{vehicleRepo: mockVehicleRepository},
			want: VehicleService{vehicleRepo: mockVehicleRepository},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewVehicleService(tt.args.vehicleRepo)
			assert.Equal(t, tt.want, got)
		},
		)
	}
}

func TestVehicleService_GetVehicleByFipeCode(t *testing.T) {
	type fields struct {
		vehicleRepo domain.VehicleRepository
	}
	type args struct {
		fipeCode   int
		orderBy    []domain.OrderBy
		pagination domain.Pagination
	}

	type TestCases struct {
		name    string
		fields  fields
		args    args
		want    []domain.Vehicle
		wantErr *errs.AppError
	}
	var tests []TestCases

	domainVehicleExamples := domain.GetDomainVehiclesExamples()

	mockVehicleRepository, ctrl := getMockVehicleRepository(t)
	t.Cleanup(ctrl.Finish)
	// -------------------------------------------------
	mockVehicleRepository.EXPECT().
		GetVehicle(
			[]domain.Condition{{Column: "fipe_code", Value: 1}},
			[]domain.OrderBy{},
			domain.Pagination{},
		).
		Return([]domain.Vehicle{domainVehicleExamples[0]}, nil).
		Times(1)
	tests = append(tests, TestCases{
		name:   "Return 1 vehicles",
		fields: fields{vehicleRepo: mockVehicleRepository},
		args: args{
			fipeCode:   1,
			orderBy:    []domain.OrderBy{},
			pagination: domain.Pagination{},
		},
		want:    []domain.Vehicle{domainVehicleExamples[0]},
		wantErr: nil,
	})
	// -------------------------------------------------
	mockVehicleRepository.EXPECT().
		GetVehicle(
			[]domain.Condition{{Column: "fipe_code", Value: 2}},
			[]domain.OrderBy{},
			domain.Pagination{},
		).
		Return([]domain.Vehicle{domainVehicleExamples[1], domainVehicleExamples[2]}, nil).
		Times(1)
	tests = append(tests, TestCases{
		name:   "Return 2 vehicles",
		fields: fields{vehicleRepo: mockVehicleRepository},
		args: args{
			fipeCode:   2,
			orderBy:    []domain.OrderBy{},
			pagination: domain.Pagination{},
		},
		want:    []domain.Vehicle{domainVehicleExamples[1], domainVehicleExamples[2]},
		wantErr: nil,
	})

	mockVehicleRepository.EXPECT().
		GetVehicle(
			[]domain.Condition{{Column: "fipe_code", Value: 3}},
			[]domain.OrderBy{},
			domain.Pagination{},
		).
		Return(nil, errs.NewNotFoundError("Vehicles not found for that fipe code")).
		Times(1)

	tests = append(tests, TestCases{
		name:   "Return 0 vehicles, NotFoundError",
		fields: fields{vehicleRepo: mockVehicleRepository},
		args: args{
			fipeCode:   3,
			orderBy:    []domain.OrderBy{},
			pagination: domain.Pagination{},
		},
		want:    nil,
		wantErr: errs.NewNotFoundError("Vehicles not found for that fipe code"),
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := VehicleService{
				vehicleRepo: tt.fields.vehicleRepo,
			}
			got, err := v.GetVehicleByFipeCode(tt.args.fipeCode, tt.args.orderBy, tt.args.pagination)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestVehicleService_GetVehicleByReferenceMonth(t *testing.T) {
	type fields struct {
		vehicleRepo domain.VehicleRepository
	}
	type args struct {
		year       int
		month      int
		orderBy    []domain.OrderBy
		pagination domain.Pagination
	}

	type TestCases struct {
		name    string
		fields  fields
		args    args
		want    []domain.Vehicle
		wantErr *errs.AppError
	}
	var tests []TestCases

	domainVehicleExamples := domain.GetDomainVehiclesExamples()

	mockVehicleRepository, ctrl := getMockVehicleRepository(t)
	t.Cleanup(ctrl.Finish)

	mockVehicleRepository.EXPECT().
		GetVehicle(
			[]domain.Condition{
				{Column: "year", Value: 2021},
				{Column: "month", Value: 7},
			},
			[]domain.OrderBy{},
			domain.Pagination{},
		).
		Return([]domain.Vehicle{domainVehicleExamples[2]}, nil).
		Times(1)
	tests = append(tests, TestCases{
		name:   "Return 1 vehicles",
		fields: fields{vehicleRepo: mockVehicleRepository},
		args: args{
			year:       2021,
			month:      7,
			orderBy:    []domain.OrderBy{},
			pagination: domain.Pagination{},
		},
		want:    []domain.Vehicle{domainVehicleExamples[2]},
		wantErr: nil,
	})

	mockVehicleRepository.EXPECT().
		GetVehicle(
			[]domain.Condition{
				{Column: "year", Value: 2021},
				{Column: "month", Value: 6},
			},
			[]domain.OrderBy{},
			domain.Pagination{},
		).
		Return([]domain.Vehicle{domainVehicleExamples[0], domainVehicleExamples[1]}, nil).
		Times(1)
	tests = append(tests, TestCases{
		name:   "Return 2 vehicles",
		fields: fields{vehicleRepo: mockVehicleRepository},
		args: args{
			year:       2021,
			month:      6,
			orderBy:    []domain.OrderBy{},
			pagination: domain.Pagination{},
		},
		want:    []domain.Vehicle{domainVehicleExamples[0], domainVehicleExamples[1]},
		wantErr: nil,
	})

	mockVehicleRepository.EXPECT().
		GetVehicle(
			[]domain.Condition{
				{Column: "year", Value: 2025},
				{Column: "month", Value: 6},
			},
			[]domain.OrderBy{},
			domain.Pagination{},
		).
		Return(nil, errs.NewNotFoundError("Vehicles not found for that reference month")).
		Times(1)

	tests = append(tests, TestCases{
		name:   "Return 0 vehicles, NotFoundError",
		fields: fields{vehicleRepo: mockVehicleRepository},
		args: args{
			year:       2025,
			month:      6,
			orderBy:    []domain.OrderBy{},
			pagination: domain.Pagination{},
		},
		want:    nil,
		wantErr: errs.NewNotFoundError("Vehicles not found for that reference month"),
	})

	mockVehicleRepository.EXPECT().
		GetVehicle(
			[]domain.Condition{
				{Column: "year", Value: 2021},
				{Column: "month", Value: 8},
			},
			[]domain.OrderBy{},
			domain.Pagination{},
		).
		Return(nil, errs.NewNotFoundError("Vehicles not found for that reference month")).
		Times(1)

	tests = append(tests, TestCases{
		name:   "Return 0 vehicles, NotFoundError",
		fields: fields{vehicleRepo: mockVehicleRepository},
		args: args{
			year:       2021,
			month:      8,
			orderBy:    []domain.OrderBy{},
			pagination: domain.Pagination{},
		},
		want:    nil,
		wantErr: errs.NewNotFoundError("Vehicles not found for that reference month"),
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := VehicleService{vehicleRepo: tt.fields.vehicleRepo}
			got, err := v.GetVehicleByReferenceYearMonth(tt.args.year, tt.args.month, tt.args.orderBy, tt.args.pagination)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
