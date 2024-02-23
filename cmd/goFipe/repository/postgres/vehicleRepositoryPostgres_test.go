package postgres

import (
	"fmt"
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
			want: errs.NewUnprocessableEntityError(fmt.Sprintf("invalid limit. The limit must be between 1 and %d", MaxLimit)),
		},
		{
			name: "Zero Limit",
			args: args{
				pagination: domain.Pagination{Offset: 0, Limit: 0},
			},
			want: errs.NewUnprocessableEntityError(fmt.Sprintf("invalid limit. The limit must be between 1 and %d", MaxLimit)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, validatePagination(tt.args.pagination), "validatePagination(%v)", tt.args.pagination)
		})
	}
}
