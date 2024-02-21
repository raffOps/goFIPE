package service

import (
	"github.com/golang/mock/gomock"
	"github.com/raffops/gofipe/cmd/goFipe/domain"
	"github.com/raffops/gofipe/cmd/goFipe/errs"
	mock_port "github.com/raffops/gofipe/cmd/goFipe/mocks"
	"github.com/raffops/gofipe/cmd/goFipe/port"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getMockVehicleRepository(t *testing.T) (*mock_port.MockVehicleRepository, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockVehicleRepository := mock_port.NewMockVehicleRepository(ctrl)
	return mockVehicleRepository, ctrl
}

func TestNewVehicleService(t *testing.T) {
	mockVehicleRepository, ctrl := getMockVehicleRepository(t)
	t.Cleanup(ctrl.Finish)

	type args struct {
		vehicleRepo port.VehicleRepository
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

func TestVehicleService_GetVehicle(t *testing.T) {
	type fields struct {
		vehicleRepo func(repo *mock_port.MockVehicleRepository)
	}
	type args struct {
		where   map[string]string
		orderBy map[string]bool
		offset  int
		limit   int
	}

	domainVehicleExamples := domain.GetDomainVehiclesExamples()

	type TestCases struct {
		name    string
		fields  fields
		args    args
		want    []domain.Vehicle
		wantErr *errs.AppError
	}
	tests := []TestCases{
		{
			name: "where fipe = 1, return 1 vehicles",
			fields: fields{
				vehicleRepo: func(repo *mock_port.MockVehicleRepository) {
					repo.EXPECT().GetVehicle(
						[]domain.WhereClause{
							{Column: "fipe_code", Value: "1", Operator: "="},
						},
						nil,
						domain.Pagination{
							Offset: 0,
							Limit:  100,
						},
					).
						Return(
							[]domain.Vehicle{domainVehicleExamples[0]},
							nil,
						).Times(1)
				}},
			args: args{
				where: map[string]string{"fipe_code": "1"},
			},
			want:    []domain.Vehicle{domainVehicleExamples[0]},
			wantErr: nil,
		},
		{
			name: "where fipe = 1, return 2 vehicles",
			fields: fields{
				vehicleRepo: func(repo *mock_port.MockVehicleRepository) {
					repo.EXPECT().GetVehicle(
						[]domain.WhereClause{
							{Column: "fipe_code", Operator: "=", Value: "2"},
						},
						nil,
						domain.Pagination{
							Offset: 0,
							Limit:  100,
						},
					).
						Return(
							[]domain.Vehicle{domainVehicleExamples[1], domainVehicleExamples[2]},
							nil,
						).Times(1)
				},
			},
			args: args{
				where: map[string]string{"fipe_code": "2"},
			},
			want:    []domain.Vehicle{domainVehicleExamples[1], domainVehicleExamples[2]},
			wantErr: nil,
		},
		{
			name: "where fipe = 999, return 0 vehicles, NotFoundError",
			fields: fields{
				vehicleRepo: func(repo *mock_port.MockVehicleRepository) {
					repo.EXPECT().
						GetVehicle(
							[]domain.WhereClause{
								{
									Column:   "fipe_code",
									Operator: "=",
									Value:    "999",
								},
							},
							nil,
							domain.Pagination{
								Offset: 0,
								Limit:  100,
							},
						).Return(
						nil,
						errs.NewNotFoundError("Vehicles not found for fipe_code equal to 999"),
					).Times(1)
				},
			},
			args: args{
				where: map[string]string{"fipe_code": "999"},
			},
			want:    nil,
			wantErr: errs.NewNotFoundError("Vehicles not found for fipe_code equal to 999"),
		},
		{
			name: "where fipe = 1, pagination 10, return 2 vehicles",
			fields: fields{
				vehicleRepo: func(repo *mock_port.MockVehicleRepository) {
					repo.EXPECT().GetVehicle(
						[]domain.WhereClause{
							{Column: "fipe_code", Value: "3", Operator: "="},
						},
						nil,
						domain.Pagination{
							Offset: 0,
							Limit:  10,
						},
					).
						Return(
							[]domain.Vehicle{domainVehicleExamples[3]},
							nil,
						).Times(1)
				}},
			args: args{
				where: map[string]string{"fipe_code": "3"},
				limit: 10,
			},
			want:    []domain.Vehicle{domainVehicleExamples[3]},
			wantErr: nil,
		},
		{
			name: "where year = 2021 month 7, return 2 vehicles",
			fields: fields{
				vehicleRepo: func(repo *mock_port.MockVehicleRepository) {
					repo.EXPECT().GetVehicle(
						[]domain.WhereClause{
							{Column: "year", Operator: "=", Value: "2021"},
							{Column: "month", Operator: "=", Value: "7"},
						},
						nil,
						domain.Pagination{
							Offset: 0,
							Limit:  100,
						},
					).
						Return(
							[]domain.Vehicle{domainVehicleExamples[2], domainVehicleExamples[0]},
							nil,
						).Times(1)
				},
			},
			args: args{
				where: map[string]string{"year": "2021", "month": "7"},
			},
			want:    []domain.Vehicle{domainVehicleExamples[2], domainVehicleExamples[0]},
			wantErr: nil,
		},
		{
			name: "where year = 2021 month 7, order by mean value desc, return 2 vehicles",
			fields: fields{
				vehicleRepo: func(repo *mock_port.MockVehicleRepository) {
					repo.EXPECT().GetVehicle(
						[]domain.WhereClause{
							{Column: "year", Value: "2021", Operator: "="},
							{Column: "month", Value: "7", Operator: "="},
						},
						[]domain.OrderByClause{
							{Column: "mean_value", IsDesc: true},
						},
						domain.Pagination{
							Offset: 0,
							Limit:  100,
						},
					).
						Return(
							[]domain.Vehicle{domainVehicleExamples[2], domainVehicleExamples[0]},
							nil,
						).Times(1)
				},
			},
			args: args{
				where:   map[string]string{"year": "2021", "month": "7"},
				orderBy: map[string]bool{"mean_value": true},
			},
			want:    []domain.Vehicle{domainVehicleExamples[2], domainVehicleExamples[0]},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		mockVehicleRepository, ctrl := getMockVehicleRepository(t)
		t.Cleanup(ctrl.Finish)
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.vehicleRepo(mockVehicleRepository)
			v := VehicleService{vehicleRepo: mockVehicleRepository}
			got, err := v.GetVehicle(tt.args.where, tt.args.orderBy, tt.args.offset, tt.args.limit)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

//
//	mockVehicleRepository.EXPECT().
//		GetVehicle(
//			[]domain.WhereClause{
//				{Column: "year", Value: 2021, Operator: "="},
//				{Column: "month", Value: 7, Operator: "="},
//			},
//			[]domain.OrderByClause{},
//			domain.Pagination{},
//		).
//		Return([]domain.Vehicle{domainVehicleExamples[2]}, nil).
//		Times(1)
//	tests = append(tests, TestCases{
//		name:   "Return 1 vehicles",
//		fields: fields{vehicleRepo: mockVehicleRepository},
//		args: args{
//			year:       2021,
//			month:      7,
//			orderBy:    []domain.OrderByClause{},
//			pagination: domain.Pagination{},
//		},
//		want:    []domain.Vehicle{domainVehicleExamples[2]},
//		wantErr: nil,
//	})
//	// -------------------------------------------------
//	mockVehicleRepository.EXPECT().
//		GetVehicle(
//			[]domain.WhereClause{
//				{Column: "year", Value: 2021, Operator: "="},
//				{Column: "month", Value: 6, Operator: "="},
//			},
//			[]domain.OrderByClause{},
//			domain.Pagination{},
//		).
//		Return([]domain.Vehicle{domainVehicleExamples[0], domainVehicleExamples[1]}, nil).
//		Times(1)
//	tests = append(tests, TestCases{
//		name:   "Return 2 vehicles",
//		fields: fields{vehicleRepo: mockVehicleRepository},
//		args: args{
//			year:       2021,
//			month:      6,
//			orderBy:    []domain.OrderByClause{},
//			pagination: domain.Pagination{},
//		},
//		want:    []domain.Vehicle{domainVehicleExamples[0], domainVehicleExamples[1]},
//		wantErr: nil,
//	})
//	// -------------------------------------------------
//	mockVehicleRepository.EXPECT().
//		GetVehicle(
//			[]domain.WhereClause{
//				{Column: "year", Value: 2025, Operator: "="},
//				{Column: "month", Value: 6, Operator: "="},
//			},
//			[]domain.OrderByClause{},
//			domain.Pagination{},
//		).
//		Return(nil, errs.NewNotFoundError("Vehicles not found for that reference month")).
//		Times(1)
//
//	tests = append(tests, TestCases{
//		name:   "Return 0 vehicles, NotFoundError",
//		fields: fields{vehicleRepo: mockVehicleRepository},
//		args: args{
//			year:       2025,
//			month:      6,
//			orderBy:    []domain.OrderByClause{},
//			pagination: domain.Pagination{},
//		},
//		want:    nil,
//		wantErr: errs.NewNotFoundError("Vehicles not found for that reference month"),
//	})
//	// -------------------------------------------------
//	mockVehicleRepository.EXPECT().
//		GetVehicle(
//			[]domain.WhereClause{
//				{Column: "year", Value: 2021, Operator: "="},
//				{Column: "month", Value: 8, Operator: "="},
//			},
//			[]domain.OrderByClause{},
//			domain.Pagination{},
//		).
//		Return(nil, errs.NewNotFoundError("Vehicles not found for that reference month")).
//		Times(1)
//
//	tests = append(tests, TestCases{
//		name:   "Return 0 vehicles, NotFoundError",
//		fields: fields{vehicleRepo: mockVehicleRepository},
//		args: args{
//			year:       2021,
//			month:      8,
//			orderBy:    []domain.OrderByClause{},
//			pagination: domain.Pagination{},
//		},
//		want:    nil,
//		wantErr: errs.NewNotFoundError("Vehicles not found for that reference month"),
//	})
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			v := VehicleService{vehicleRepo: tt.fields.vehicleRepo}
//			got, err := v.GetVehicleByReferenceYearMonth(tt.args.year, tt.args.month, tt.args.orderBy, tt.args.pagination)
//			assert.Equal(t, tt.want, got)
//			assert.Equal(t, tt.wantErr, err)
//		})
//	}
//}
