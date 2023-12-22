package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rjribeiro/goFIPE/cmd/gofipe/database"
	"github.com/rjribeiro/goFIPE/cmd/gofipe/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func init() {
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_USER", "root")
	os.Setenv("POSTGRES_PASSWORD", "root")
	os.Setenv("POSTGRES_DB", "gofipe_test")
	database.ConnectToDB()
	gin.SetMode(gin.TestMode)
}

func TestGetHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/health-check", GetHealthCheck)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/health-check", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"api up\"}", w.Body.String())
}

func TestCreateVehicle(t *testing.T) {
	router := gin.Default()
	router.POST("/veiculo", CreateVehicle)

	var testCases = []struct {
		input      []byte
		wantedCode int
	}{
		{
			input: []byte(`{
				"ano":               2023,
				"mes":               10,
				"valor":             80877.0,
				"marca":             "Renault",
				"modelo":            "STEPWAY Intense Flex 1.6 16V  Aut.",
				"ano_modelo":        "2022",
				"combustivel":       "Gasolina",
				"codigo_fipe":       "025282-4",
				"mes_referencia":    202311,
				"tipo_veiculo":      1,
				"sigla_combustivel": "G",
				"data_consulta":     "2023-12-17 17:23:17"
			}`),
			wantedCode: http.StatusCreated,
		},
		{
			input: []byte(`{
				"ano":               2023,
				"mes":               10,
				"valor":             80877,
				"marca":             "foo",
				"modelo":            "STEPWAY Intense Flex 1.6 16V  Aut.",
				"ano_modelo":        "2022",
				"combustivel":       "Gasolina",
				"codigo_fipe":       "025282-4",
				"mes_referencia":    202311,
				"tipo_veiculo":      1,
				"sigla_combustivel": "G",
				"data_consulta":     "2023-12-17 17:23:17"
			}`),
			wantedCode: http.StatusCreated,
		},
		{
			input: []byte(`{
				"ano":               "foo",
				"mes":               10,
				"valor":             80877.0,
				"marca":             "Renault",
				"modelo":            "STEPWAY Intense Flex 1.6 16V  Aut.",
				"ano_modelo":        "2022",
				"combustivel":       "Gasolina",
				"codigo_fipe":       "025282-4",
				"mes_referencia":    "bar"
				"tipo_veiculo":      1,
				"sigla_combustivel": "G",
				"data_consulta":     "2023-12-17 17:23:17"
			}`),
			wantedCode: http.StatusBadRequest,
		},

		{
			input: []byte(`{
				"ano":               "foo",
				"mes":               10,
				"valor":             80877.0,
				"marca":             "Renault",
				"modelo":            "STEPWAY Intense Flex 1.6 16V  Aut.",
				"ano_modelo":        "2022",
				"combustivel":       "Gasolina",
				"codigo_fipe":       "025282-4",
				"mes_referencia":    "bar"
				"tipo_veiculo":      1,
				"sigla_combustivel": "G",
				"data_consulta":     "2023-12-17 17:23:17"
			}`),
			wantedCode: http.StatusBadRequest,
		},
		{
			input: []byte(`{
				"ano":               2023,
				"mes":               10,
				"valor":             true,
				"marca":             "Renault",
				"modelo":            "STEPWAY Intense Flex 1.6 16V  Aut.",
				"ano_modelo":        "2022",
				"combustivel":       "Gasolina",
				"codigo_fipe":       "025282-4",
				"mes_referencia":    202311,
				"tipo_veiculo":      1,
				"sigla_combustivel": "G",
				"data_consulta":     "2023-12-17 17:23:17"
			}`),
			wantedCode: http.StatusBadRequest,
		},
	}

	for index, testCase := range testCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/veiculo", bytes.NewReader(testCase.input))
		router.ServeHTTP(w, req)

		assert.Equal(t, testCase.wantedCode, w.Code, "Test %d", index)
		if testCase.wantedCode == 201 {
			got := models.Vehicle{}
			defer database.DB.Unscoped().Delete(&got)

			unmarshelledInput := models.Vehicle{}
			json.Unmarshal(testCase.input, &unmarshelledInput)
			database.DB.FirstOrCreate(&got)
			assert.Equal(t, got.AnoModelo, unmarshelledInput.AnoModelo)
			assert.Equal(t, got.Ano, unmarshelledInput.Ano)
			assert.Equal(t, got.Mes, unmarshelledInput.Mes)
			assert.Equal(t, got.Valor, unmarshelledInput.Valor)
			assert.Equal(t, got.Marca, unmarshelledInput.Marca)
			database.DB.Unscoped().Delete(&got)
		}
	}
}
