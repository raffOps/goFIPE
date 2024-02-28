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
	mockPort "github.com/raffops/gofipe/cmd/goFipe/mocks"
	"github.com/raffops/gofipe/cmd/goFipe/port"
	"github.com/raffops/gofipe/cmd/goFipe/utils"
	"github.com/stretchr/testify/assert"
)

func getMockVehicleService(t *testing.T) (*mockPort.MockVehicleService, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockVehicleService := mockPort.NewMockVehicleService(ctrl)
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
			if got := NewVehicleHandler(tt.args.vehicleService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVehicleHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVehicleHandler_Get(t *testing.T) {
	type Dependencies struct {
		vehicleService func(service *mockPort.MockVehicleService)
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
				"where":  "fipe_code:111111-1",
				"order":  "year:asc",
				"offset": "0",
				"limit":  "1",
			},
			dependencies: Dependencies{
				vehicleService: func(service *mockPort.MockVehicleService) {
					service.EXPECT().
						GetVehicle(
							map[string]string{
								"fipe_code": "111111-1",
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
			wantBody:       "[{\"ano\":2021,\"mes\":7,\"fipe_code\":\"111111-1\",\"marca\":\"Acura\",\"modelo\":\"Integra GS 1.8\",\"ano_modelo\":\"1992 Gasolina\",\"autenticacao\":\"1\",\"valor_medio\":700}]\n",
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
				vehicleService: func(service *mockPort.MockVehicleService) {
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
			wantBody:       "[{\"ano\":2021,\"mes\":7,\"fipe_code\":\"111111-1\",\"marca\":\"Acura\",\"modelo\":\"Integra GS 1.8\",\"ano_modelo\":\"1992 Gasolina\",\"autenticacao\":\"1\",\"valor_medio\":700}]\n",
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
				vehicleService: func(service *mockPort.MockVehicleService) {},
			},
			wantBody:       "Clausula where 0 deve ser no formato 'key:value'\n",
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
				vehicleService: func(service *mockPort.MockVehicleService) {},
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
				vehicleService: func(service *mockPort.MockVehicleService) {},
			},
			wantBody:       "Clausula order 0: Value deve ser asc ou desc\n",
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
				vehicleService: func(service *mockPort.MockVehicleService) {},
			},
			wantBody:       "Clausula order 1: Value deve ser asc ou desc\n",
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
				vehicleService: func(service *mockPort.MockVehicleService) {},
			},
			wantBody:       "Clausula order 1 deve ser no formato 'key:value'\n",
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
				vehicleService: func(service *mockPort.MockVehicleService) {},
			},
			wantBody:       "Campo order deve possuir no minimo 1 clausula\n",
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
				vehicleService: func(service *mockPort.MockVehicleService) {},
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
				vehicleService: func(service *mockPort.MockVehicleService) {},
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
				vehicleService: func(service *mockPort.MockVehicleService) {
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
	appHost := os.Getenv("APP_HOST")
	appPort := os.Getenv("APP_PORT")
	if appHost == "" || appPort == "" {
		t.Fatal("APP_HOST or APP_PORT environment variable is not set")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVehicleService, ctrl := getMockVehicleService(t)
			t.Cleanup(ctrl.Finish)
			t.Logf(appHost, appPort)
			path := fmt.Sprintf("%s:%s/vehicles", appHost, appPort)
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
		name    string
		input   string
		want    map[string]string
		wantErr *errs.AppError
	}{
		{
			name:  "Normal use case",
			input: "fipe_code:1,year:2021,month:7",
			want: map[string]string{
				"fipe_code": "1",
				"year":      "2021",
				"month":     "7",
			},
			wantErr: nil,
		},
		{
			name:    "Invalid where clause",
			input:   "fipe_code:1,year:2021,month",
			want:    nil,
			wantErr: errs.NewBadRequestError("Clausula where 2 deve ser no formato 'key:value'"),
		},
		{
			name:    "Empty where clause",
			input:   "",
			want:    nil,
			wantErr: errs.NewBadRequestError("Campo where deve possuir no minimo 1 clausula"),
		},
		{
			name:    "Invalid where clause",
			input:   "fipe_code:1,year:2021,month:",
			want:    nil,
			wantErr: errs.NewBadRequestError("Clausula where 2 deve ser no formato 'key:value'"),
		},
		{
			name:    "Invalid where clause",
			input:   "fipe_code:1,year:2021,month:7,",
			want:    nil,
			wantErr: errs.NewBadRequestError("Clausula where 3 deve ser no formato 'key:value'"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := handleWhereParameter(tt.input)
			assert.Equalf(t, tt.want, got, "handleWhereParameter(%v)", tt.input)
			assert.Equalf(t, tt.wantErr, gotErr, "handleWhereParameter(%v)", tt.input)
		})
	}
}

func Test_handleOrderByParameter(t *testing.T) {
	type args struct {
		orderByString string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]bool
		wantErr *errs.AppError
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
			wantErr: nil,
		},
		{
			name: "Test with abnormal case, clauses with invalid format",
			args: args{
				orderByString: "key1:asc,key2:desc,key3",
			},
			want:    nil,
			wantErr: errs.NewBadRequestError("Clausula order 2 deve ser no formato 'key:value'"),
		},
		{
			name: "Test with abnormal case, empty value",
			args: args{
				orderByString: "key1:asc,key2:",
			},
			want:    nil,
			wantErr: errs.NewBadRequestError("Clausula order 1 deve ser no formato 'key:value'"),
		},
		{
			name: "Test without clauses",
			args: args{
				orderByString: "",
			},
			want:    nil,
			wantErr: errs.NewBadRequestError("Campo order deve possuir no minimo 1 clausula"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := handleOrderByParameter(tt.args.orderByString)
			assert.Equalf(t, tt.want, got, "handleOrderByParameter(%v)", tt.args.orderByString)
			assert.Equalf(t, tt.wantErr, got1, "handleOrderByParameter(%v)", tt.args.orderByString)
		})
	}
}
