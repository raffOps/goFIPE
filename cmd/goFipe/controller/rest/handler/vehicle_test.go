package handler

import (
	"fmt"
	"github.com/raffops/gofipe/cmd/goFipe/errs"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/raffops/gofipe/cmd/goFipe/domain"
	mock_port "github.com/raffops/gofipe/cmd/goFipe/mocks"
	"github.com/raffops/gofipe/cmd/goFipe/port"
	"github.com/raffops/gofipe/cmd/goFipe/utils"
	"github.com/stretchr/testify/assert"
)

func getMockVehicleService(t *testing.T) (*mock_port.MockVehicleService, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockVehicleService := mock_port.NewMockVehicleService(ctrl)
	return mockVehicleService, ctrl
}

func TestNewHandler(t *testing.T) {
	type args struct {
		vehicleService port.VehicleService
	}

	mockVehicleService, ctrl := getMockVehicleService(t)
	t.Cleanup(ctrl.Finish)

	tests := []struct {
		name string
		args args
		want VehicleHandler
	}{
		{
			name: "Single test",
			args: args{vehicleService: mockVehicleService},
			want: VehicleHandler{vehicleService: mockVehicleService},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHandler(tt.args.vehicleService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVehicleHandler_Get(t *testing.T) {
	type Dependencies struct {
		vehicleService func(service *mock_port.MockVehicleService)
	}

	vehiclesExamples := domain.GetDomainVehiclesExamples()

	tests := []struct {
		name           string
		args           map[string]interface{}
		dependencies   Dependencies
		wantBody       string
		wantStatusCode int
	}{
		{
			name: "where by fipe_code",
			args: map[string]interface{}{
				"where":  "fipe_code:1",
				"order":  "year:asc",
				"offset": "0",
				"limit":  "1",
			},
			dependencies: Dependencies{
				vehicleService: func(service *mock_port.MockVehicleService) {
					service.EXPECT().
						GetVehicle(
							map[string]string{
								"fipe_code": "1",
							},
							map[string]bool{"year": false},
							0,
							1,
						).Return(
						[]domain.Vehicle{vehiclesExamples[0]},
						nil,
					)
				},
			},
			wantBody:       "[{\"ano\":2021,\"mes\":7,\"fipe_code\":1,\"marca\":\"Acura\",\"modelo\":\"Integra GS 1.8\",\"ano_modelo\":\"1992 Gasolina\",\"autenticacao\":\"1\",\"valor_medio\":700}]\n",
			wantStatusCode: http.StatusOK,
		},
		{
			name: "where by year and month",
			args: map[string]interface{}{
				"where":  "year:2021,month:7",
				"order":  "year:asc,month:asc",
				"offset": "0",
				"limit":  "1",
			},
			dependencies: Dependencies{
				vehicleService: func(service *mock_port.MockVehicleService) {
					service.EXPECT().
						GetVehicle(
							map[string]string{
								"year":  "2021",
								"month": "7",
							},
							map[string]bool{"year": false, "month": false},
							0,
							1,
						).Return(
						[]domain.Vehicle{vehiclesExamples[0]},
						nil,
					)
				},
			},
			wantBody:       "[{\"ano\":2021,\"mes\":7,\"fipe_code\":1,\"marca\":\"Acura\",\"modelo\":\"Integra GS 1.8\",\"ano_modelo\":\"1992 Gasolina\",\"autenticacao\":\"1\",\"valor_medio\":700}]\n",
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Invalid where clause",
			args: map[string]interface{}{
				"where":  "fipe_code:",
				"order":  "year:asc",
				"offset": "0",
				"limit":  "1",
			},
			dependencies: Dependencies{
				vehicleService: func(service *mock_port.MockVehicleService) {},
			},
			wantBody:       "Clausula Where invalida. fipe_code:. Clausula where deve ser no formato 'key:value'\n",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid where by clause",
			args: map[string]interface{}{
				"where":  "",
				"order":  "year:desc",
				"offset": "0",
				"limit":  "1",
			},
			dependencies: Dependencies{
				vehicleService: func(service *mock_port.MockVehicleService) {},
			},
			wantBody:       "Campo where deve possuir no minimo 1 clausula\n",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid order by clause",
			args: map[string]interface{}{
				"where":  "fipe_code:1",
				"order":  "year:invalid",
				"offset": "0",
				"limit":  "1",
			},
			dependencies: Dependencies{
				vehicleService: func(service *mock_port.MockVehicleService) {},
			},
			wantBody:       "Clausula OrderBy invalida. year:invalid. Value deve ser asc ou desc\n",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid order by clause",
			args: map[string]interface{}{
				"where":  "fipe_code:1",
				"order":  "year:asc,month:invalid",
				"offset": "0",
				"limit":  "1",
			},
			dependencies: Dependencies{
				vehicleService: func(service *mock_port.MockVehicleService) {},
			},
			wantBody:       "Clausula OrderBy invalida. month:invalid. Value deve ser asc ou desc\n",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid order by clause",
			args: map[string]interface{}{
				"where":  "fipe_code:1",
				"order":  "year:asc,month:",
				"offset": "0",
				"limit":  "1",
			},
			dependencies: Dependencies{
				vehicleService: func(service *mock_port.MockVehicleService) {},
			},
			wantBody:       "Clausula OrderBy invalida. month:. Clausula where deve ser no formato 'key:value'\n",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid order by clause",
			args: map[string]interface{}{
				"where":  "fipe_code:1",
				"order":  "",
				"offset": "0",
				"limit":  "1",
			},
			dependencies: Dependencies{
				vehicleService: func(service *mock_port.MockVehicleService) {},
			},
			wantBody:       "Campo OrderBy deve possuir no minimo 1 clausula\n",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid offset",
			args: map[string]interface{}{
				"where":  "fipe_code:1",
				"order":  "year:asc",
				"offset": "invalid",
				"limit":  "1",
			},
			dependencies: Dependencies{
				vehicleService: func(service *mock_port.MockVehicleService) {},
			},
			wantBody:       "Offset deve ser um numero inteiro\n",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid limit",
			args: map[string]interface{}{
				"where":  "fipe_code:1",
				"order":  "year:asc",
				"offset": "0",
				"limit":  "invalid",
			},
			dependencies: Dependencies{
				vehicleService: func(service *mock_port.MockVehicleService) {},
			},
			wantBody:       "Limit deve ser um numero inteiro\n",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Unexpected error",
			args: map[string]interface{}{
				"where":  "fipe_code:1",
				"order":  "year:asc",
				"offset": "0",
				"limit":  "1",
			},
			dependencies: Dependencies{
				vehicleService: func(service *mock_port.MockVehicleService) {
					service.EXPECT().
						GetVehicle(
							map[string]string{
								"fipe_code": "1",
							},
							map[string]bool{"year": false},
							0,
							1,
						).
						Return(
							nil,
							errs.NewUnexpectedError("Unexpected error"),
						)
				},
			},
			wantBody:       "Unexpected error\n",
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVehicleService, ctrl := getMockVehicleService(t)
			t.Cleanup(ctrl.Finish)
			path := fmt.Sprintf("%s:%s/vehicles", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
			urlEncoded := utils.EncodeUrl(path, tt.args)
			req, err := http.NewRequest("GET", urlEncoded, nil)
			if err != nil {
				t.Fatal(t)
			}
			rr := httptest.NewRecorder()
			tt.dependencies.vehicleService(mockVehicleService)
			vehicleHandler := VehicleHandler{vehicleService: mockVehicleService}
			handler := http.HandlerFunc(vehicleHandler.Get)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatusCode, rr.Code)
			assert.Equal(t, tt.wantBody, rr.Body.String())
		})
	}
}

func Test_handleWhereParameter(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected map[string]string
		isError  bool
	}{
		{
			"Test normal use case",
			"key1:value1,key2:value2",
			map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			false,
		},
		{
			"Test with abnormal case, clauses with invalid format",
			"key1 value1,key2:value2",
			nil,
			true,
		},
		{
			"Test with abnormal case, empty clause",
			"key1:,key2:value2",
			nil,
			true,
		},
		{
			"Test without clauses",
			"",
			nil,
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := handleWhereParameter(tc.input)

			if (err != nil) != tc.isError {
				t.Errorf("Expected error = %v, got %v", tc.isError, err)
			}

			if len(got) != len(tc.expected) {
				t.Errorf("Expected = %v, got %v", tc.expected, got)
			}

			for k, v := range tc.expected {
				if value, ok := got[k]; ok {
					if value != v {
						t.Errorf("Expected = %v, got %v", tc.expected, got)
					}
				} else {
					t.Errorf("Expected = %v, got %v", tc.expected, got)
				}
			}
		})
	}
}

func Test_handleOrderByParameter(t *testing.T) {
	type args struct {
		orderByString string
	}
	tests := []struct {
		name  string
		args  args
		want  map[string]bool
		want1 *errs.AppError
	}{
		{
			name: "Test normal use case",
			args: args{
				orderByString: "key1:asc,key2:desc",
			},
			want: map[string]bool{
				"key1": false,
				"key2": true,
			},
			want1: nil,
		},
		{
			name: "Test with abnormal case, invalid value",
			args: args{
				orderByString: "key1:asc,key2:invalid",
			},
			want:  nil,
			want1: errs.NewBadRequestError("Clausula OrderBy invalida. key2:invalid. Value deve ser asc ou desc"),
		},
		{
			name: "Test with abnormal case, invalid format",
			args: args{
				orderByString: "key1:asc,key2",
			},
			want:  nil,
			want1: errs.NewBadRequestError("Clausula OrderBy invalida. key2. Clausula where deve ser no formato 'key:value'"),
		},
		{
			name: "Test with abnormal case, empty clause",
			args: args{
				orderByString: "key1:asc,key2:",
			},
			want:  nil,
			want1: errs.NewBadRequestError("Clausula OrderBy invalida. key2:. Clausula where deve ser no formato 'key:value'"),
		},
		{
			name: "Test without clauses",
			args: args{
				orderByString: "",
			},
			want:  nil,
			want1: errs.NewBadRequestError("Campo OrderBy deve possuir no minimo 1 clausula"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := handleOrderByParameter(tt.args.orderByString)
			assert.Equalf(t, tt.want, got, "handleOrderByParameter(%v)", tt.args.orderByString)
			assert.Equalf(t, tt.want1, got1, "handleOrderByParameter(%v)", tt.args.orderByString)
		})
	}
}
