package domain

import "testing"

func TestIsValidFipeCode(t *testing.T) {
	type args struct {
		fipeCode string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid fipe code",
			args: args{fipeCode: "123456-7"},
			want: true,
		},
		{
			name: "Invalid fipe code (missing hyphen)",
			args: args{fipeCode: "1234567"},
			want: false,
		},
		{
			name: "Invalid fipe code (missing digit after hyphen)",
			args: args{fipeCode: "123456-"},
			want: false,
		},
		{
			name: "Invalid fipe code (missing digit before hyphen)",
			args: args{fipeCode: "-7"},
			want: false,
		},
		{
			name: "Invalid fipe code (missing digits before and after hyphen)",
			args: args{fipeCode: "-"},
			want: false,
		},
		{
			name: "Invalid fipe code (too many digits before hyphen)",
			args: args{fipeCode: "1234567-7"},
			want: false,
		},
		{
			name: "Invalid fipe code (too many digits after hyphen)",
			args: args{fipeCode: "123456-77"},
			want: false,
		},
		{
			name: "Invalid fipe code (too many digits before and after hyphen)",
			args: args{fipeCode: "1234567-77"},
			want: false,
		},
		{
			name: "Invalid fipe code (non-digit before hyphen)",
			args: args{fipeCode: "123456a-7"},
			want: false,
		},
		{
			name: "Invalid fipe code (non-digit after hyphen)",
			args: args{fipeCode: "123456-a"},
			want: false,
		},
		{
			name: "Invalid fipe code (non-digit before and after hyphen)",
			args: args{fipeCode: "123456-a7"},
			want: false,
		},
		{
			name: "Invalid fipe code (non-hyphen character)",
			args: args{fipeCode: "123456a7"},
			want: false,
		},
		{
			name: "Invalid fipe code (empty string)",
			args: args{fipeCode: ""},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidFipeCode(tt.args.fipeCode); got != tt.want {
				t.Errorf("IsValidFipeCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidYearMonth(t *testing.T) {
	type args struct {
		year  int
		month int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid year and month",
			args: args{year: 2021, month: 7},
			want: true,
		},
		{
			name: "Invalid year (future year)",
			args: args{year: 2042, month: 7},
			want: false,
		},
		{
			name: "Invalid year (future month)",
			args: args{year: 2024, month: 8},
			want: false,
		},
		{
			name: "Invalid year (before 1900)",
			args: args{year: 1899, month: 7},
			want: false,
		},
		{
			name: "Invalid month (less than 1)",
			args: args{year: 2021, month: 0},
			want: false,
		},
		{
			name: "Invalid month (greater than 12)",
			args: args{year: 2021, month: 13},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidYearMonth(tt.args.year, tt.args.month); got != tt.want {
				t.Errorf("IsValidYearMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}
