// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository.go

// Package repo is a generated GoMock package.
package repo

import (
	context "context"
	reflect "reflect"

	payments "golangreferenceapi/internal/payments"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CreatePaymentPlan mocks base method.
func (m *MockRepository) CreatePaymentPlan(ctx context.Context, arg *payments.CreatePlanParams) (*payments.Plan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePaymentPlan", ctx, arg)
	ret0, _ := ret[0].(*payments.Plan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePaymentPlan indicates an expected call of CreatePaymentPlan.
func (mr *MockRepositoryMockRecorder) CreatePaymentPlan(ctx, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePaymentPlan", reflect.TypeOf((*MockRepository)(nil).CreatePaymentPlan), ctx, arg)
}

// ListPaymentPlansByUserID mocks base method.
func (m *MockRepository) ListPaymentPlansByUserID(ctx context.Context, userID string) ([]*payments.Plan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPaymentPlansByUserID", ctx, userID)
	ret0, _ := ret[0].([]*payments.Plan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPaymentPlansByUserID indicates an expected call of ListPaymentPlansByUserID.
func (mr *MockRepositoryMockRecorder) ListPaymentPlansByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPaymentPlansByUserID", reflect.TypeOf((*MockRepository)(nil).ListPaymentPlansByUserID), ctx, userID)
}
