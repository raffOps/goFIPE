// Code generated by MockGen. DO NOT EDIT.
// Source: port.go

// Package mock_port is a generated GoMock package.
package mock_port

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/raffops/gofipe/cmd/goFipe/domain"
	errs "github.com/raffops/gofipe/cmd/goFipe/errs"
)

// MockVehicleService is a mock of VehicleService interface.
type MockVehicleService struct {
	ctrl     *gomock.Controller
	recorder *MockVehicleServiceMockRecorder
}

// MockVehicleServiceMockRecorder is the mock recorder for MockVehicleService.
type MockVehicleServiceMockRecorder struct {
	mock *MockVehicleService
}

// NewMockVehicleService creates a new mock instance.
func NewMockVehicleService(ctrl *gomock.Controller) *MockVehicleService {
	mock := &MockVehicleService{ctrl: ctrl}
	mock.recorder = &MockVehicleServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVehicleService) EXPECT() *MockVehicleServiceMockRecorder {
	return m.recorder
}

// GetVehicle mocks base method.
func (m *MockVehicleService) GetVehicle(where map[string]string, orderBy map[string]bool, limit, offset int) ([]domain.Vehicle, *errs.AppError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVehicle", where, orderBy, limit, offset)
	ret0, _ := ret[0].([]domain.Vehicle)
	ret1, _ := ret[1].(*errs.AppError)
	return ret0, ret1
}

// GetVehicle indicates an expected call of GetVehicle.
func (mr *MockVehicleServiceMockRecorder) GetVehicle(where, orderBy, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVehicle", reflect.TypeOf((*MockVehicleService)(nil).GetVehicle), where, orderBy, limit, offset)
}

// MockVehicleRepository is a mock of VehicleRepository interface.
type MockVehicleRepository struct {
	ctrl     *gomock.Controller
	recorder *MockVehicleRepositoryMockRecorder
}

// MockVehicleRepositoryMockRecorder is the mock recorder for MockVehicleRepository.
type MockVehicleRepositoryMockRecorder struct {
	mock *MockVehicleRepository
}

// NewMockVehicleRepository creates a new mock instance.
func NewMockVehicleRepository(ctrl *gomock.Controller) *MockVehicleRepository {
	mock := &MockVehicleRepository{ctrl: ctrl}
	mock.recorder = &MockVehicleRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVehicleRepository) EXPECT() *MockVehicleRepositoryMockRecorder {
	return m.recorder
}

// GetVehicle mocks base method.
func (m *MockVehicleRepository) GetVehicle(conditions []domain.WhereClause, orderBy []domain.OrderByClause, pagination domain.Pagination) ([]domain.Vehicle, *errs.AppError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVehicle", conditions, orderBy, pagination)
	ret0, _ := ret[0].([]domain.Vehicle)
	ret1, _ := ret[1].(*errs.AppError)
	return ret0, ret1
}

// GetVehicle indicates an expected call of GetVehicle.
func (mr *MockVehicleRepositoryMockRecorder) GetVehicle(conditions, orderBy, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVehicle", reflect.TypeOf((*MockVehicleRepository)(nil).GetVehicle), conditions, orderBy, pagination)
}
