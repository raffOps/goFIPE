package handler

import (
	"encoding/json"
	"fmt"
	"github.com/raffops/gofipe/cmd/goFipe/controller/rest/dto"
	"net/http"
	"strconv"
	"strings"

	"github.com/raffops/gofipe/cmd/goFipe/errs"
	"github.com/raffops/gofipe/cmd/goFipe/port"
)

type VehicleHandler struct {
	vehicleService port.VehicleService
}

func NewHandler(vehicleService port.VehicleService) VehicleHandler {
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

	var responseVehicles []dto.VehicleResponse
	for _, vehicle := range vehicles {
		responseVehicles = append(responseVehicles, dto.VehicleResponseFromDomain(vehicle))
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseVehicles)
}

func handleWhereParameter(whereString string) (map[string]string, *errs.AppError) {
	if len(strings.TrimSpace(whereString)) == 0 {
		return nil, errs.NewBadRequestError("Campo where deve possuir no minimo 1 clausula")
	}

	where := map[string]string{}
	for _, split := range strings.Split(whereString, ",") {
		split = strings.TrimSpace(split)
		if len(split) == 0 {
			return nil, errs.NewBadRequestError(fmt.Sprintf("Clausula Where invalida. %s. Clausula where deve ser no formato 'key:value'", split))
		}

		keyValue := strings.Split(split, ":")
		if len(keyValue) != 2 {
			return nil, errs.NewBadRequestError(
				fmt.Sprintf("Clausula Where invalida. %s. Clausula where deve ser no formato 'key:value'", split),
			)
		}

		if keyValue[0] == "" || keyValue[1] == "" {
			return nil, errs.NewBadRequestError(
				fmt.Sprintf("Clausula Where invalida. %s. Clausula where deve ser no formato 'key:value'", split),
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
		return nil, errs.NewBadRequestError("Campo OrderBy deve possuir no minimo 1 clausula")
	}

	orderBy := map[string]bool{}
	for _, split := range strings.Split(orderByString, ",") {
		split = strings.TrimSpace(split)
		if len(split) == 0 {
			return nil, errs.NewBadRequestError(fmt.Sprintf("Clausula OrderBy invalida. %s. Clausula where deve ser no formato 'key:value'", split))
		}

		keyValue := strings.Split(split, ":")
		if len(keyValue) != 2 {
			return nil, errs.NewBadRequestError(
				fmt.Sprintf("Clausula OrderBy invalida. %s. Clausula where deve ser no formato 'key:value'", split),
			)
		}

		if keyValue[0] == "" || keyValue[1] == "" {
			return nil, errs.NewBadRequestError(
				fmt.Sprintf("Clausula OrderBy invalida. %s. Clausula where deve ser no formato 'key:value'", split),
			)
		}

		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])

		if value != "asc" && value != "desc" {
			return nil, errs.NewBadRequestError(
				fmt.Sprintf("Clausula OrderBy invalida. %s. Value deve ser asc ou desc", split),
			)
		}

		if value == "desc" {
			orderBy[key] = true
		} else if value == "asc" {
			orderBy[key] = false
		} else {
			return nil, errs.NewBadRequestError(
				fmt.Sprintf("Clausula OrderBy invalida. %s. Value deve ser asc ou desc", split),
			)
		}
	}

	return orderBy, nil
}
