package models

import "gorm.io/gorm"

type Veiculo struct {
	gorm.Model
	Ano              uint16
	Mes              uint8
	Valor            float32
	Marca            string
	Modelo           string
	AnoModelo        string
	Combustivel      string
	CodigoFipe       string
	MesReferencia    string
	TipoVeiculo      int
	SiglaCombustivel string
	DataConsulta     string
}
