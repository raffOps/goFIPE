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

func TestCreateVeiculo(t *testing.T) {

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/veiculo", CreateVeiculo)
	w := httptest.NewRecorder()

	var testCases = []models.Veiculo{
		{
			Ano:              2023,
			Mes:              10,
			Valor:            80877.0,
			Marca:            "Renault",
			Modelo:           "STEPWAY Intense Flex 1.6 16V  Aut.",
			AnoModelo:        "2022",
			Combustivel:      "Gasolina",
			CodigoFipe:       "025282-4",
			MesReferencia:    202311,
			TipoVeiculo:      1,
			SiglaCombustivel: "G",
			DataConsulta:     "2023-12-17 17:23:17",
		},
	}

	for _, testCase := range testCases {
		marshalled, _ := json.Marshal(testCase)
		var got models.Veiculo

		req, _ := http.NewRequest("POST", "/veiculo", bytes.NewReader(marshalled))
		router.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
		database.DB.FirstOrCreate(&got)
		assert.Equal(t, got.Ano, testCase.Ano)
		assert.Equal(t, got.Mes, testCase.Mes)
		assert.Equal(t, got.Valor, testCase.Valor)
		assert.Equal(t, got.Marca, testCase.Marca)
		assert.Equal(t, got.Modelo, testCase.Modelo)

		database.DB.Unscoped().Delete(&got)
	}

}
