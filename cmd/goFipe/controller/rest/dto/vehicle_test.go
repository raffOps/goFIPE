package dto

import (
	"github.com/raffops/gofipe/cmd/goFipe/domain"
	"reflect"
	"testing"
)

func Test_vehicleResponseFromDomain(t *testing.T) {
	type args struct {
		vehicle domain.Vehicle
	}
	tests := []struct {
		name string
		args args
		want VehicleResponse
	}{
		{
			name: "Single test",
			args: args{vehicle: domain.Vehicle{
				Year:           2021,
				Month:          7,
				FipeCode:       1,
				Brand:          "Acura",
				Model:          "Integra GS 1.8",
				YearModel:      "1992 Gasolina",
				Authentication: "1",
				MeanValue:      700,
			}},
			want: VehicleResponse{
				Year:           2021,
				Month:          7,
				FipeCode:       1,
				Brand:          "Acura",
				Model:          "Integra GS 1.8",
				YearModel:      "1992 Gasolina",
				Authentication: "1",
				MeanValue:      700,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VehicleResponseFromDomain(tt.args.vehicle); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VehicleResponseFromDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
