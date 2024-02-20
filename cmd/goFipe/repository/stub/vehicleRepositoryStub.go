package stub

//
//import (
//	"github.com/raffops/gofipe/cmd/goFipe/domain"
//	"github.com/raffops/gofipe/cmd/goFipe/errs"
//	"reflect"
//	"sort"
//)
//
//type VehicleRepositoryStub struct {
//	Vehicles []domain.Vehicle
//}
//
//func NewVehicleRepositoryStub() VehicleRepositoryStub {
//	vehicles := domain.GetDomainVehiclesExamples()
//	return VehicleRepositoryStub{vehicles}
//}
//
//func (v VehicleRepositoryStub) GetVehicle(
//	conditions []domain.Condition,
//	orderBy []domain.OrderBy,
//	pagination domain.Pagination) ([]Vehicle, *errs.AppError) {
//
//	vehicles := v.where(conditions)
//
//	if len(vehicles) == 0 {
//		return nil, errs.NewNotFoundError("No vehicles found")
//	}
//
//	vehicles = v.orderBy(vehicles, orderBy)
//
//	return nil, nil
//
//}
//
//func  (v VehicleRepositoryStub) where (conditions []domain.Condition) []domain.Vehicle {
//	var vehicles []domain.Vehicle
//
//	for _, vehicle := range v.Vehicles {
//		for _, condition := range conditions {
//			if reflect.ValueOf(vehicle).
//				FieldByName(condition.Column).Interface() != condition.Value {
//					break
//			}
//			vehicles = append(vehicles, vehicle)
//		}
//	}
//	return vehicles
//}
//
//func (v VehicleRepositoryStub) orderBy(vehicles []domain.Vehicle, by []domain.OrderBy) []domain.Vehicle {
//	for _, order := range by {
//		if order.Order == "desc" {
//			vehicles = 	sort.Slice(vehicles, func(i, j int) bool {
//				return vehicles[i]. < employees[j].Age
//			})
//
//			vehicles = v.orderByAsc(vehicles, order.Column)
//		}
//	}
//	return vehicles
//}
//
//
