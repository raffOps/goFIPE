package models

import "gorm.io/gorm"

type Veiculo struct {
	gorm.Model
	Ano              uint16  `json:"ano"`
	Mes              uint8   `json:"mes"`
	Valor            float32 `json:"valor"`
	Marca            string  `json:"marca"`
	Modelo           string  `json:"modelo"`
	AnoModelo        string  `json:"ano_modelo"`
	Combustivel      string  `json:"combustivel"`
	CodigoFipe       string  `json:"codigo_fipe"`
	MesReferencia    int     `json:"mes_referencia"`
	TipoVeiculo      int     `json:"tipo_veiculo"`
	SiglaCombustivel string  `json:"sigla_combustivel"`
	DataConsulta     string  `json:"data_consulta"`
}
