// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/zsoltggs/tabeo-interview/services/bookings/internal/transport/v1/healthhttp (interfaces: HealthCheckable)
//
// Generated by this command:
//
//	mockgen -package=mocks -destination=../../../mocks/healthhttp.go github.com/zsoltggs/tabeo-interview/services/bookings/internal/transport/v1/healthhttp HealthCheckable
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockHealthCheckable is a mock of HealthCheckable interface.
type MockHealthCheckable struct {
	ctrl     *gomock.Controller
	recorder *MockHealthCheckableMockRecorder
}

// MockHealthCheckableMockRecorder is the mock recorder for MockHealthCheckable.
type MockHealthCheckableMockRecorder struct {
	mock *MockHealthCheckable
}

// NewMockHealthCheckable creates a new mock instance.
func NewMockHealthCheckable(ctrl *gomock.Controller) *MockHealthCheckable {
	mock := &MockHealthCheckable{ctrl: ctrl}
	mock.recorder = &MockHealthCheckableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHealthCheckable) EXPECT() *MockHealthCheckableMockRecorder {
	return m.recorder
}

// Health mocks base method.
func (m *MockHealthCheckable) Health() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Health")
	ret0, _ := ret[0].(error)
	return ret0
}

// Health indicates an expected call of Health.
func (mr *MockHealthCheckableMockRecorder) Health() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Health", reflect.TypeOf((*MockHealthCheckable)(nil).Health))
}
