package stub

//
//import (
//	"github.com/raffops/gofipe/cmd/goFipe/domain"
//	"github.com/raffops/gofipe/cmd/goFipe/errs"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//func TestNewVehicleRepositoryStub(t *testing.T) {
//	tests := []struct {
//		name string
//		want VehicleRepositoryStub
//	}{
//
//		{
//			name: "Return VehicleRepositoryStub",
//			want: VehicleRepositoryStub{
//				Vehicles: domain.GetDomainVehiclesExamples(),
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got := NewVehicleRepositoryStub()
//			assert.Equal(t, tt.want, got)
//		})
//	}
//}
//
////func TestVehicleRepositoryStub_GetVehicleByFipeCode(t *testing.T) {
////	type fields struct {
////		Vehicles []domain.Vehicle
////	}
////	type args struct {
////		fipeCode int
////	}
////	tests := []struct {
////		name    string
////		fields  fields
////		args    args
////		want    []domain.Vehicle
////		wantErr error
////	}{
////		{
////			name: "Return 1 vehicle",
////			fields: fields{
////				Vehicles: vehiclesStub,
////			},
////			args: args{fipeCode: 2},
////			want: []domain.Vehicle{
////				{
////					ReferenceMonth: 202106,
////					FipeCode:       2,
////				},
////			},
////			wantErr: nil,
////		},
////		{
////			name: "Return 2 vehicles",
////			fields: fields{
////				Vehicles: vehiclesStub,
////			},
////			args: args{fipeCode: 1},
////			want: []domain.Vehicle{
////				{
////					ReferenceMonth: 202106,
////					FipeCode:       1,
////				},
////				{
////					ReferenceMonth: 202107,
////					FipeCode:       1,
////				},
////			},
////			wantErr: nil,
////		},
////		{
////			name: "Return 0 vehicles",
////			fields: fields{
////				Vehicles: vehiclesStub,
////			},
////			args:    args{fipeCode: 3},
////			want:    nil,
////			wantErr: nil,
////		},
////	}
////
////	for _, tt := range tests {
////		t.Run(tt.name, func(t *testing.T) {
////			v := VehicleRepositoryStub{
////				Vehicles: tt.fields.Vehicles,
////			}
////			got, err := v.GetVehicleByFipeCode(tt.args.fipeCode)
////			assert.Equal(t, tt.want, got)
////			assert.Equal(t, tt.wantErr, err)
////		})
////	}
////}
////
////func TestVehicleRepositoryStub_GetVehicleByReferenceMonth(t *testing.T) {
////	type fields struct {
////		Vehicles []domain.Vehicle
////	}
////	type args struct {
////		referenceMonth int
////	}
////	tests := []struct {
////		name    string
////		fields  fields
////		args    args
////		want    []domain.Vehicle
////		wantErr error
////	}{
////		{
////			name: "Return 1 vehicle",
////			fields: fields{
////				Vehicles: vehiclesStub,
////			},
////			args: args{referenceMonth: 202107},
////			want: []domain.Vehicle{
////				{
////					ReferenceMonth: 202107,
////					FipeCode:       1,
////				},
////			},
////			wantErr: nil,
////		},
////		{
////			name: "Return 2 vehicles",
////			fields: fields{
////				Vehicles: vehiclesStub,
////			},
////			args: args{referenceMonth: 202106},
////			want: []domain.Vehicle{
////				{
////					ReferenceMonth: 202106,
////					FipeCode:       1,
////				},
////				{
////					ReferenceMonth: 202106,
////					FipeCode:       2,
////				},
////			},
////			wantErr: nil,
////		},
////	}
////	for _, tt := range tests {
////		t.Run(tt.name, func(t *testing.T) {
////			v := VehicleRepositoryStub{
////				Vehicles: tt.fields.Vehicles,
////			}
////			got, err := v.GetVehicleByReferenceMonth(tt.args.referenceMonth)
////			assert.Equal(t, tt.want, got)
////			assert.Equal(t, tt.wantErr, err)
////		})
////	}
////}
//
//func TestVehicleRepositoryStub_GetVehicle(t *testing.T) {
//	type fields struct {
//		Vehicles []domain.Vehicle
//	}
//	type args struct {
//		conditions []domain.Condition
//		orderBy    []domain.OrderBy
//		pagination domain.Pagination
//	}
//	tests := []struct {
//		name      string
//		fields    fields
//		args      args
//		want      []domain.Vehicle
//		wantError *errs.AppError
//	}{
//		{
//			name: "Return 1 vehicle by fipe code",
//			fields: fields{
//				Vehicles: domain.GetDomainVehiclesExamples(),
//			},
//			args: args{
//				conditions: []domain.Condition{{Column: "fipe_code", Value: 1}},
//				orderBy:    []domain.OrderBy{},
//				pagination: domain.Pagination{},
//			},
//			want: []domain.Vehicle{
//				{
//					Year:           2021,
//					Month:          7,
//					FipeCode:       1,
//					Brand:          "Acura",
//					Model:          "Integra GS 1.8",
//					YearModel:      "1992 Gasolina",
//					Authentication: "1",
//					MeanValue:      700,
//				},
//			},
//			wantError: nil,
//		},
//		{
//			name: "Return 2 vehicles by fipe code order by year and month desc",
//			fields: fields{
//				Vehicles: domain.GetDomainVehiclesExamples(),
//			},
//			args: args{
//				conditions: []domain.Condition{{Column: "fipe_code", Value: 2}},
//				orderBy:    []domain.OrderBy{{Column: "year", Order: "desc"}, {Column: "month", Order: "desc"}},
//				pagination: domain.Pagination{},
//			},
//			want: []domain.Vehicle{
//				{
//					Year:           2021,
//					Month:          8,
//					FipeCode:       2,
//					Brand:          "Fiat",
//					Model:          "147 C/ CL",
//					YearModel:      "1991 Gasolina",
//					Authentication: "2",
//					MeanValue:      800,
//				},
//				{
//					Year:           2021,
//					Month:          7,
//					FipeCode:       2,
//					Brand:          "Fiat",
//					Model:          "147 C/ CL",
//					YearModel:      "1991 Gasolina",
//					Authentication: "2",
//					MeanValue:      801,
//				},
//				{
//					Year:           2021,
//					Month:          6,
//					FipeCode:       2,
//					Brand:          "Fiat",
//					Model:          "147 C/ CL",
//					YearModel:      "1991 Gasolina",
//					Authentication: "2",
//					MeanValue:      800,
//				},
//			},
//			wantError: nil,
//		},
//		{
//			name: "Return 2 vehicles by fipe code order by year and month desc with pagination",
//			fields: fields{
//				Vehicles: domain.GetDomainVehiclesExamples(),
//			},
//			args: args{
//				conditions: []domain.Condition{{Column: "fipe_code", Value: 2}},
//				orderBy:    []domain.OrderBy{{Column: "year", Order: "desc"}, {Column: "month", Order: "desc"}},
//				pagination: domain.Pagination{Limit: 1, Offset: 1},
//			},
//			want: []domain.Vehicle{
//				{
//					Year:           2021,
//					Month:          7,
//					FipeCode:       2,
//					Brand:          "Fiat",
//					Model:          "147 C/ CL",
//					YearModel:      "1991 Gasolina",
//					Authentication: "2",
//					MeanValue:      801,
//				},
//			},
//			wantError: nil,
//		},
//		{
//			name: "Return 2 vehicles by invalid columns",
//			fields: fields{
//				Vehicles: domain.GetDomainVehiclesExamples(),
//			},
//			args: args{
//				conditions: []domain.Condition{{Column: "invalid_column", Value: 2}},
//			},
//			want:      nil,
//			wantError: errs.NewInvalidColumnError("invalid column: invalid_column"),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			v := VehicleRepositoryStub{
//				Vehicles: tt.fields.Vehicles,
//			}
//			got, err := v.GetVehicle(tt.args.conditions, tt.args.orderBy, tt.args.pagination)
//			assert.Equalf(t, tt.want, got, "GetVehicle(%v, %v, %v)", tt.args.conditions, tt.args.orderBy, tt.args.pagination)
//			assert.Equalf(t, tt.wantError, err, "GetVehicle(%v, %v, %v)", tt.args.conditions, tt.args.orderBy, tt.args.pagination)
//		})
//	}
//}
