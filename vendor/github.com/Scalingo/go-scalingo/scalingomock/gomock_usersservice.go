// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Scalingo/go-scalingo (interfaces: UsersService)

// Package scalingomock is a generated GoMock package.
package scalingomock

import (
	go_scalingo "github.com/Scalingo/go-scalingo"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockUsersService is a mock of UsersService interface
type MockUsersService struct {
	ctrl     *gomock.Controller
	recorder *MockUsersServiceMockRecorder
}

// MockUsersServiceMockRecorder is the mock recorder for MockUsersService
type MockUsersServiceMockRecorder struct {
	mock *MockUsersService
}

// NewMockUsersService creates a new mock instance
func NewMockUsersService(ctrl *gomock.Controller) *MockUsersService {
	mock := &MockUsersService{ctrl: ctrl}
	mock.recorder = &MockUsersServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUsersService) EXPECT() *MockUsersServiceMockRecorder {
	return m.recorder
}

// Self mocks base method
func (m *MockUsersService) Self() (*go_scalingo.User, error) {
	ret := m.ctrl.Call(m, "Self")
	ret0, _ := ret[0].(*go_scalingo.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Self indicates an expected call of Self
func (mr *MockUsersServiceMockRecorder) Self() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Self", reflect.TypeOf((*MockUsersService)(nil).Self))
}

// UpdateUser mocks base method
func (m *MockUsersService) UpdateUser(arg0 go_scalingo.UpdateUserParams) (*go_scalingo.User, error) {
	ret := m.ctrl.Call(m, "UpdateUser", arg0)
	ret0, _ := ret[0].(*go_scalingo.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser
func (mr *MockUsersServiceMockRecorder) UpdateUser(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUsersService)(nil).UpdateUser), arg0)
}
