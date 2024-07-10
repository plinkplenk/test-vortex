// Code generated by MockGen. DO NOT EDIT.
// Source: orders.go
//
// Generated by this command:
//
//	mockgen -source=orders.go -destination=mocks/mock.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	orders "github.com/plinkplenk/test-vortex/internal/orders"
	gomock "go.uber.org/mock/gomock"
)

// MockOrdersService is a mock of OrdersService interface.
type MockOrdersService struct {
	ctrl     *gomock.Controller
	recorder *MockOrdersServiceMockRecorder
}

// MockOrdersServiceMockRecorder is the mock recorder for MockOrdersService.
type MockOrdersServiceMockRecorder struct {
	mock *MockOrdersService
}

// NewMockOrdersService creates a new mock instance.
func NewMockOrdersService(ctrl *gomock.Controller) *MockOrdersService {
	mock := &MockOrdersService{ctrl: ctrl}
	mock.recorder = &MockOrdersServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrdersService) EXPECT() *MockOrdersServiceMockRecorder {
	return m.recorder
}

// GetOrderBook mocks base method.
func (m *MockOrdersService) GetOrderBook(ctx context.Context, exchangeName, pair string) ([]orders.Depth, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderBook", ctx, exchangeName, pair)
	ret0, _ := ret[0].([]orders.Depth)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderBook indicates an expected call of GetOrderBook.
func (mr *MockOrdersServiceMockRecorder) GetOrderBook(ctx, exchangeName, pair any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderBook", reflect.TypeOf((*MockOrdersService)(nil).GetOrderBook), ctx, exchangeName, pair)
}

// GetOrderHistory mocks base method.
func (m *MockOrdersService) GetOrderHistory(ctx context.Context, client orders.Client) ([]*orders.History, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderHistory", ctx, client)
	ret0, _ := ret[0].([]*orders.History)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderHistory indicates an expected call of GetOrderHistory.
func (mr *MockOrdersServiceMockRecorder) GetOrderHistory(ctx, client any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderHistory", reflect.TypeOf((*MockOrdersService)(nil).GetOrderHistory), ctx, client)
}

// SaveOrder mocks base method.
func (m *MockOrdersService) SaveOrder(ctx context.Context, client orders.Client, order *orders.History) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveOrder", ctx, client, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveOrder indicates an expected call of SaveOrder.
func (mr *MockOrdersServiceMockRecorder) SaveOrder(ctx, client, order any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveOrder", reflect.TypeOf((*MockOrdersService)(nil).SaveOrder), ctx, client, order)
}

// SaveOrderBook mocks base method.
func (m *MockOrdersService) SaveOrderBook(ctx context.Context, exchangeName, pair string, orderBook []orders.Depth) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveOrderBook", ctx, exchangeName, pair, orderBook)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveOrderBook indicates an expected call of SaveOrderBook.
func (mr *MockOrdersServiceMockRecorder) SaveOrderBook(ctx, exchangeName, pair, orderBook any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveOrderBook", reflect.TypeOf((*MockOrdersService)(nil).SaveOrderBook), ctx, exchangeName, pair, orderBook)
}