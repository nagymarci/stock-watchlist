// Code generated by MockGen. DO NOT EDIT.
// Source: service/notifier.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	model "github.com/nagymarci/stock-user-profile/model"
	model0 "github.com/nagymarci/stock-watchlist/model"
	logrus "github.com/sirupsen/logrus"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	reflect "reflect"
)

// MockwatchlistList is a mock of watchlistList interface
type MockwatchlistList struct {
	ctrl     *gomock.Controller
	recorder *MockwatchlistListMockRecorder
}

// MockwatchlistListMockRecorder is the mock recorder for MockwatchlistList
type MockwatchlistListMockRecorder struct {
	mock *MockwatchlistList
}

// NewMockwatchlistList creates a new mock instance
func NewMockwatchlistList(ctrl *gomock.Controller) *MockwatchlistList {
	mock := &MockwatchlistList{ctrl: ctrl}
	mock.recorder = &MockwatchlistListMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockwatchlistList) EXPECT() *MockwatchlistListMockRecorder {
	return m.recorder
}

// List mocks base method
func (m *MockwatchlistList) List() ([]model0.Watchlist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]model0.Watchlist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockwatchlistListMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockwatchlistList)(nil).List))
}

// MockrecommendationProvider is a mock of recommendationProvider interface
type MockrecommendationProvider struct {
	ctrl     *gomock.Controller
	recorder *MockrecommendationProviderMockRecorder
}

// MockrecommendationProviderMockRecorder is the mock recorder for MockrecommendationProvider
type MockrecommendationProviderMockRecorder struct {
	mock *MockrecommendationProvider
}

// NewMockrecommendationProvider creates a new mock instance
func NewMockrecommendationProvider(ctrl *gomock.Controller) *MockrecommendationProvider {
	mock := &MockrecommendationProvider{ctrl: ctrl}
	mock.recorder = &MockrecommendationProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockrecommendationProvider) EXPECT() *MockrecommendationProviderMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockrecommendationProvider) Get(id primitive.ObjectID) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockrecommendationProviderMockRecorder) Get(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockrecommendationProvider)(nil).Get), id)
}

// Update mocks base method
func (m *MockrecommendationProvider) Update(log *logrus.Entry, id primitive.ObjectID, stocks []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", log, id, stocks)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockrecommendationProviderMockRecorder) Update(log, id, stocks interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockrecommendationProvider)(nil).Update), log, id, stocks)
}

// MockemailSender is a mock of emailSender interface
type MockemailSender struct {
	ctrl     *gomock.Controller
	recorder *MockemailSenderMockRecorder
}

// MockemailSenderMockRecorder is the mock recorder for MockemailSender
type MockemailSenderMockRecorder struct {
	mock *MockemailSender
}

// NewMockemailSender creates a new mock instance
func NewMockemailSender(ctrl *gomock.Controller) *MockemailSender {
	mock := &MockemailSender{ctrl: ctrl}
	mock.recorder = &MockemailSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockemailSender) EXPECT() *MockemailSenderMockRecorder {
	return m.recorder
}

// SendNotification mocks base method
func (m *MockemailSender) SendNotification(profileName string, removed, added, currentStocks []string, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendNotification", profileName, removed, added, currentStocks, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendNotification indicates an expected call of SendNotification
func (mr *MockemailSenderMockRecorder) SendNotification(profileName, removed, added, currentStocks, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendNotification", reflect.TypeOf((*MockemailSender)(nil).SendNotification), profileName, removed, added, currentStocks, email)
}

// MockstockGetter is a mock of stockGetter interface
type MockstockGetter struct {
	ctrl     *gomock.Controller
	recorder *MockstockGetterMockRecorder
}

// MockstockGetterMockRecorder is the mock recorder for MockstockGetter
type MockstockGetterMockRecorder struct {
	mock *MockstockGetter
}

// NewMockstockGetter creates a new mock instance
func NewMockstockGetter(ctrl *gomock.Controller) *MockstockGetter {
	mock := &MockstockGetter{ctrl: ctrl}
	mock.recorder = &MockstockGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockstockGetter) EXPECT() *MockstockGetterMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockstockGetter) Get(symbol string) (model0.StockData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", symbol)
	ret0, _ := ret[0].(model0.StockData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockstockGetterMockRecorder) Get(symbol interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockstockGetter)(nil).Get), symbol)
}

// MockstockRecommendator is a mock of stockRecommendator interface
type MockstockRecommendator struct {
	ctrl     *gomock.Controller
	recorder *MockstockRecommendatorMockRecorder
}

// MockstockRecommendatorMockRecorder is the mock recorder for MockstockRecommendator
type MockstockRecommendatorMockRecorder struct {
	mock *MockstockRecommendator
}

// NewMockstockRecommendator creates a new mock instance
func NewMockstockRecommendator(ctrl *gomock.Controller) *MockstockRecommendator {
	mock := &MockstockRecommendator{ctrl: ctrl}
	mock.recorder = &MockstockRecommendatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockstockRecommendator) EXPECT() *MockstockRecommendatorMockRecorder {
	return m.recorder
}

// GetAllRecommendedStock mocks base method
func (m *MockstockRecommendator) GetAllRecommendedStock(stocks []model0.StockData, numReqs int, userprofile *model.Userprofile) []model0.CalculatedStockInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllRecommendedStock", stocks, numReqs, userprofile)
	ret0, _ := ret[0].([]model0.CalculatedStockInfo)
	return ret0
}

// GetAllRecommendedStock indicates an expected call of GetAllRecommendedStock
func (mr *MockstockRecommendatorMockRecorder) GetAllRecommendedStock(stocks, numReqs, userprofile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllRecommendedStock", reflect.TypeOf((*MockstockRecommendator)(nil).GetAllRecommendedStock), stocks, numReqs, userprofile)
}

// MockuserprofileGetter is a mock of userprofileGetter interface
type MockuserprofileGetter struct {
	ctrl     *gomock.Controller
	recorder *MockuserprofileGetterMockRecorder
}

// MockuserprofileGetterMockRecorder is the mock recorder for MockuserprofileGetter
type MockuserprofileGetterMockRecorder struct {
	mock *MockuserprofileGetter
}

// NewMockuserprofileGetter creates a new mock instance
func NewMockuserprofileGetter(ctrl *gomock.Controller) *MockuserprofileGetter {
	mock := &MockuserprofileGetter{ctrl: ctrl}
	mock.recorder = &MockuserprofileGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockuserprofileGetter) EXPECT() *MockuserprofileGetterMockRecorder {
	return m.recorder
}

// GetUserprofile mocks base method
func (m *MockuserprofileGetter) GetUserprofile(userId string) (model.Userprofile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserprofile", userId)
	ret0, _ := ret[0].(model.Userprofile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserprofile indicates an expected call of GetUserprofile
func (mr *MockuserprofileGetterMockRecorder) GetUserprofile(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserprofile", reflect.TypeOf((*MockuserprofileGetter)(nil).GetUserprofile), userId)
}
