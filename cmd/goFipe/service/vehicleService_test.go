package service

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/raffops/gofipe/cmd/goFipe/domain"
	"github.com/raffops/gofipe/cmd/goFipe/errs"
	mockPort "github.com/raffops/gofipe/cmd/goFipe/mocks"
	"github.com/raffops/gofipe/cmd/goFipe/port"
	"github.com/stretchr/testify/assert"
)

func getMockVehicleRepository(t *testing.T) (*mockPort.MockVehicleRepository, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockVehicleRepository := mockPort.NewMockVehicleRepository(ctrl)
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
		vehicleRepo func(repo *mockPort.MockVehicleRepository)
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
			name: "where year = 2021 month 7, order by mean value desc, limit 0, return 2 vehicles",
			fields: fields{
				vehicleRepo: func(repo *mockPort.MockVehicleRepository) {
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
							Limit:  domain.MaxLimit - 1,
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
				offset:  0,
				limit:   domain.MaxLimit - 1,
			},
			want:    []domain.Vehicle{domainVehicleExamples[2], domainVehicleExamples[0]},
			wantErr: nil,
		},
		{
			name: "where fipe = 111111-1, return 1 vehicles",
			fields: fields{
				vehicleRepo: func(repo *mockPort.MockVehicleRepository) {
					repo.EXPECT().GetVehicle(
						[]domain.WhereClause{
							{Column: "fipe_code", Value: "111111-1", Operator: "="},
						},
						[]domain.OrderByClause{
							{Column: "mean_value", IsDesc: true},
						},
						domain.Pagination{
							Offset: 0,
							Limit:  domain.MaxLimit - 1,
						},
					).
						Return(
							[]domain.Vehicle{domainVehicleExamples[0]},
							nil,
						).Times(1)
				}},
			args: args{
				where:   map[string]string{"fipe_code": "111111-1"},
				orderBy: map[string]bool{"mean_value": true},
				offset:  0,
				limit:   domain.MaxLimit - 1,
			},
			want:    []domain.Vehicle{domainVehicleExamples[0]},
			wantErr: nil,
		},
		{
			name: "where fipe = 222222-2, return 2 vehicles",
			fields: fields{
				vehicleRepo: func(repo *mockPort.MockVehicleRepository) {
					repo.EXPECT().GetVehicle(
						[]domain.WhereClause{
							{Column: "fipe_code", Operator: "=", Value: "222222-2"},
						},
						[]domain.OrderByClause{
							{Column: "mean_value", IsDesc: true},
						},
						domain.Pagination{
							Offset: 0,
							Limit:  domain.MaxLimit - 1,
						},
					).
						Return(
							[]domain.Vehicle{domainVehicleExamples[1], domainVehicleExamples[2]},
							nil,
						).Times(1)
				},
			},
			args: args{
				where:   map[string]string{"fipe_code": "222222-2"},
				orderBy: map[string]bool{"mean_value": true},
				offset:  0,
				limit:   domain.MaxLimit - 1,
			},
			want:    []domain.Vehicle{domainVehicleExamples[1], domainVehicleExamples[2]},
			wantErr: nil,
		},
		{
			name: "where fipe = 999999-9, return 0 vehicles, NotFoundError",
			fields: fields{
				vehicleRepo: func(repo *mockPort.MockVehicleRepository) {
					repo.EXPECT().
						GetVehicle(
							[]domain.WhereClause{
								{
									Column:   "fipe_code",
									Operator: "=",
									Value:    "999999-9",
								},
							},
							[]domain.OrderByClause{
								{Column: "mean_value", IsDesc: true},
							},
							domain.Pagination{
								Offset: 0,
								Limit:  domain.MaxLimit - 1,
							},
						).Return(
						nil,
						errs.NewNotFoundError("Vehicles not found for fipe_code equal to 999999-9"),
					).Times(1)
				},
			},
			args: args{
				where:   map[string]string{"fipe_code": "999999-9"},
				orderBy: map[string]bool{"mean_value": true},
				offset:  0,
				limit:   domain.MaxLimit - 1,
			},
			want:    nil,
			wantErr: errs.NewNotFoundError("Vehicles not found for fipe_code equal to 999999-9"),
		},
		{
			name: "where fipe = 333333-3, pagination 10, return 2 vehicles",
			fields: fields{
				vehicleRepo: func(repo *mockPort.MockVehicleRepository) {
					repo.EXPECT().GetVehicle(
						[]domain.WhereClause{
							{Column: "fipe_code", Value: "333333-3", Operator: "="},
						},
						[]domain.OrderByClause{
							{Column: "mean_value", IsDesc: true},
						},
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
				where:   map[string]string{"fipe_code": "333333-3"},
				orderBy: map[string]bool{"mean_value": true},
				offset:  0,
				limit:   10,
			},
			want:    []domain.Vehicle{domainVehicleExamples[3]},
			wantErr: nil,
		},
		{
			name: "where year = 2021 month 7, return 2 vehicles",
			fields: fields{
				vehicleRepo: func(repo *mockPort.MockVehicleRepository) {
					repo.EXPECT().GetVehicle(
						[]domain.WhereClause{
							{Column: "year", Operator: "=", Value: "2021"},
							{Column: "month", Operator: "=", Value: "7"},
						},
						[]domain.OrderByClause{
							{Column: "mean_value", IsDesc: true},
						},
						domain.Pagination{
							Offset: 0,
							Limit:  domain.MaxLimit - 1,
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
				offset:  0,
				limit:   domain.MaxLimit - 1,
			},
			want:    []domain.Vehicle{domainVehicleExamples[2], domainVehicleExamples[0]},
			wantErr: nil,
		},
		{
			name: "where year = 2021 month 7, order by mean value desc, return 2 vehicles",
			fields: fields{
				vehicleRepo: func(repo *mockPort.MockVehicleRepository) {
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
							Limit:  domain.MaxLimit - 1,
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
				offset:  0,
				limit:   domain.MaxLimit - 1,
			},
			want:    []domain.Vehicle{domainVehicleExamples[2], domainVehicleExamples[0]},
			wantErr: nil,
		},
		{
			name: "Empty where, BadRequestError",
			fields: fields{
				vehicleRepo: func(repo *mockPort.MockVehicleRepository) {
					repo.EXPECT().GetVehicle(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				},
			},
			args: args{
				where:   map[string]string{},
				orderBy: map[string]bool{"mean_value": true},
				offset:  0,
				limit:   domain.MaxLimit - 1,
			},
			want:    nil,
			wantErr: errs.NewBadRequestError("Where is required"),
		},
		{
			name: "Invalid fipe code, ValidationError",
			fields: fields{
				vehicleRepo: func(repo *mockPort.MockVehicleRepository) {
					repo.EXPECT().GetVehicle(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				},
			},
			args: args{
				where:   map[string]string{"fipe_code": "invalid"},
				orderBy: map[string]bool{"mean_value": true},
				offset:  0,
				limit:   domain.MaxLimit - 1,
			},
			want:    nil,
			wantErr: errs.NewValidationError("Invalid fipe code"),
		},
		{
			name: "Invalid year, ValidationError",
			fields: fields{
				vehicleRepo: func(repo *mockPort.MockVehicleRepository) {
					repo.EXPECT().GetVehicle(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				},
			},
			args: args{
				where:   map[string]string{"year": "invalid"},
				orderBy: map[string]bool{"mean_value": true},
				offset:  0,
				limit:   domain.MaxLimit - 1,
			},
			want:    nil,
			wantErr: errs.NewValidationError("Invalid year"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVehicleRepository, ctrl := getMockVehicleRepository(t)
			t.Cleanup(ctrl.Finish)
			tt.fields.vehicleRepo(mockVehicleRepository)
			v := VehicleService{vehicleRepo: mockVehicleRepository}
			got, err := v.GetVehicle(tt.args.where, tt.args.orderBy, tt.args.offset, tt.args.limit)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_validatePagination(t *testing.T) {
	type args struct {
		offset int
		limit  int
	}
	tests := []struct {
		name string
		args args
		want *errs.AppError
	}{
		{
			name: "valid offset and limit, no error",
			args: args{
				offset: 0,
				limit:  3,
			},
			want: nil,
		},
		{
			name: "offset smaller than 0, BadRequestError",
			args: args{
				offset: -1,
				limit:  3,
			},
			want: errs.NewValidationError("Offset must be greater than 0"),
		},

		{
			name: "offset greater than limit, BadRequestError",
			args: args{
				offset: 3,
				limit:  2,
			},
			want: errs.NewValidationError("Offset must be smaller than Limit"),
		},
		{
			name: "limit greater than domain.MaxLimit, BadRequestError",
			args: args{
				offset: 0,
				limit:  domain.MaxLimit + 1,
			},
			want: errs.NewValidationError(fmt.Sprintf("Limit must be smaller or equal than %d", domain.MaxLimit)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, validatePagination(tt.args.offset, tt.args.limit), "validatePagination(%v, %v)", tt.args.offset, tt.args.limit)
		})
	}
}

func Test_validateOrderBy(t *testing.T) {
	type args struct {
		orderBy map[string]bool
	}
	tests := []struct {
		name string
		args args
		want *errs.AppError
	}{
		{
			name: "valid orderBy, no error",
			args: args{
				orderBy: map[string]bool{"fipe_code": true},
			},
			want: nil,
		},
		{
			name: "invalid column, BadRequestError",
			args: args{
				orderBy: map[string]bool{"invalid_column": true},
			},
			want: errs.NewValidationError("Invalid column: invalid_column"),
		},
		{
			name: "empty orderBy, BadRequestError",
			args: args{
				orderBy: map[string]bool{},
			},
			want: errs.NewBadRequestError("OrderBy is required"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, validateOrderBy(tt.args.orderBy), "validateOrderBy(%v)", tt.args.orderBy)
		})
	}
}

func Test_validateWhere(t *testing.T) {
	type args struct {
		where map[string]string
	}
	tests := []struct {
		name string
		args args
		want *errs.AppError
	}{
		{
			name: "valid where by fipe_code, no error",
			args: args{
				where: map[string]string{"fipe_code": "111111-1"},
			},
			want: nil,
		},
		{
			name: "valid where by year, no error",
			args: args{
				where: map[string]string{"year": "2021"},
			},
			want: nil,
		},
		{
			name: "valid where by month, no error",
			args: args{
				where: map[string]string{"month": "7"},
			},
			want: nil,
		},
		{
			name: "valid where by mean_value, no error",
			args: args{
				where: map[string]string{"mean_value": "1000.0"},
			},
			want: nil,
		},
		{
			name: "empty where, BadRequestError",
			args: args{
				where: map[string]string{},
			},
			want: errs.NewBadRequestError("Where is required"),
		},
		{
			name: "invalid fipe_code, ValidationError",
			args: args{
				where: map[string]string{"fipe_code": "invalid"},
			},
			want: errs.NewValidationError("Invalid fipe code"),
		},
		{
			name: "invalid year, ValidationError",
			args: args{
				where: map[string]string{"year": "invalid"},
			},
			want: errs.NewValidationError("Invalid year"),
		},
		{
			name: "invalid year, ValidationError",
			args: args{
				where: map[string]string{"year": "-1"},
			},
			want: errs.NewValidationError("Invalid year"),
		},
		{
			name: "invalid month, ValidationError",
			args: args{
				where: map[string]string{"month": "invalid"},
			},
			want: errs.NewValidationError("Invalid month"),
		},
		{
			name: "invalid month, ValidationError",
			args: args{
				where: map[string]string{"month": "-1"},
			},
			want: errs.NewValidationError("Invalid month"),
		},
		{
			name: "invalid month, ValidationError",
			args: args{
				where: map[string]string{"month": "13"},
			},
			want: errs.NewValidationError("Invalid month"),
		},
		{
			name: "invalid mean_value, ValidationError",
			args: args{
				where: map[string]string{"mean_value": "invalid"},
			},
			want: errs.NewValidationError("Invalid mean value"),
		},
		{
			name: "invalid column, ValidationError",
			args: args{
				where: map[string]string{"invalid_column": "invalid"},
			},
			want: errs.NewValidationError("Invalid Column"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, validateWhere(tt.args.where), "validateWhere(%v)", tt.args.where)
		})
	}
}
