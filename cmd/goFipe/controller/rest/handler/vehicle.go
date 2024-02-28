package handler

import (
	"encoding/json"
	"fmt"
	"github.com/raffops/gofipe/cmd/goFipe/controller/rest/dto"
	"github.com/raffops/gofipe/cmd/goFipe/domain/ports"
	"net/http"
	"strconv"
	"strings"

	"github.com/raffops/gofipe/cmd/goFipe/errs"
)

type VehicleHandler struct {
	vehicleService ports.VehicleService
}

func NewVehicleHandler(vehicleService ports.VehicleService) VehicleHandler {
	return VehicleHandler{vehicleService: vehicleService}
}

func (h VehicleHandler) Get(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")

	where, errWhere := handleWhereParameter(r.URL.Query().Get("where"))
	if errWhere != nil {
		http.Error(w, errWhere.Message, errWhere.Code)
		return
	}

	orderBy, errOrderBy := handleOrderByParameter(r.URL.Query().Get("order"))
	if errOrderBy != nil {
		http.Error(w, errOrderBy.Message, errOrderBy.Code)
		return
	}

	offsetString := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetString)
	if err != nil {
		http.Error(w, "Offset deve ser um numero inteiro", http.StatusBadRequest)
		return
	}
	limitString := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		http.Error(w, "Limit deve ser um numero inteiro", http.StatusBadRequest)
		return
	}

	vehicles, errGet := h.vehicleService.GetVehicle(where, orderBy, offset, limit)
	if errGet != nil {
		http.Error(w, errGet.Message, errGet.Code)
		return
	}

	var responseVehicles []dto.GetVehicleResponse
	for _, vehicle := range vehicles {
		responseVehicles = append(responseVehicles, dto.VehicleResponseFromDomain(vehicle))
	}
	err = json.NewEncoder(w).Encode(responseVehicles)
	if err != nil {
		http.Error(w, "Erro ao serializar resposta", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleWhereParameter(whereString string) (map[string]string, *errs.AppError) {
	if len(strings.TrimSpace(whereString)) == 0 {
		return nil, errs.NewBadRequestError("Campo where deve possuir no minimo 1 clausula")
	}

	where := map[string]string{}
	for index, split := range strings.Split(whereString, ",") {
		split = strings.TrimSpace(split)
		if len(split) == 0 {
			return nil, errs.NewBadRequestError(
				fmt.Sprintf("Clausula where %d deve ser no formato 'key:value'", index),
			)
		}

		keyValue := strings.Split(split, ":")
		if len(keyValue) != 2 {
			return nil, errs.NewBadRequestError(
				fmt.Sprintf("Clausula where %d deve ser no formato 'key:value'", index),
			)
		}

		if keyValue[0] == "" || keyValue[1] == "" {
			return nil, errs.NewBadRequestError(
				fmt.Sprintf("Clausula where %d deve ser no formato 'key:value'", index),
			)
		}

		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])
		where[key] = value
	}

	return where, nil
}

func handleOrderByParameter(orderByString string) (map[string]bool, *errs.AppError) {
	if len(strings.TrimSpace(orderByString)) == 0 {
		return nil, errs.NewBadRequestError("Campo order deve possuir no minimo 1 clausula")
	}

	orderBy := map[string]bool{}
	for index, split := range strings.Split(orderByString, ",") {
		split = strings.TrimSpace(split)
		if len(split) == 0 {
			return nil, errs.NewBadRequestError(
				fmt.Sprintf("Clausula order %d deve ser no formato 'key:value'", index))
		}

		keyValue := strings.Split(split, ":")
		if len(keyValue) != 2 {
			return nil, errs.NewBadRequestError(
				fmt.Sprintf("Clausula order %d deve ser no formato 'key:value'", index),
			)
		}

		if keyValue[0] == "" || keyValue[1] == "" {
			return nil, errs.NewBadRequestError(
				fmt.Sprintf("Clausula order %d deve ser no formato 'key:value'", index),
			)
		}

		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])

		switch value {
		case "asc":
			orderBy[key] = false
		case "desc":
			orderBy[key] = true
		default:
			return nil, errs.NewBadRequestError(
				fmt.Sprintf("Clausula order %d: Value deve ser asc ou desc", index),
			)
		}
	}

	return orderBy, nil
}
