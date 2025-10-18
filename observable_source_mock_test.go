package config

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

type ObservableSourceMock struct {
	ctrl     *gomock.Controller
	recorder *ObservableSourceMockRecorder
}

type ObservableSourceMockRecorder struct {
	mock *ObservableSourceMock
}

func NewObservableSourceMock(ctrl *gomock.Controller) *ObservableSourceMock {
	mock := &ObservableSourceMock{ctrl: ctrl}
	mock.recorder = &ObservableSourceMockRecorder{mock}
	return mock
}

func (m *ObservableSourceMock) EXPECT() *ObservableSourceMockRecorder {
	return m.recorder
}

func (m *ObservableSourceMock) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *ObservableSourceMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*ObservableSourceMock)(nil).Close))
}

func (m *ObservableSourceMock) Get(path string, def ...any) any {
	m.ctrl.T.Helper()
	varargs := append([]interface{}{path}, def...)
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0 := ret[0]
	return ret0
}

func (mr *ObservableSourceMockRecorder) Get(path interface{}, def ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{path}, def...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*ObservableSourceMock)(nil).Get), varargs...)
}

func (m *ObservableSourceMock) GetPriority() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPriority")
	ret0, _ := ret[0].(int)
	return ret0
}

func (mr *ObservableSourceMockRecorder) GetPriority() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPriority", reflect.TypeOf((*ObservableSourceMock)(nil).GetPriority))
}

func (m *ObservableSourceMock) Has(path string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Has", path)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (mr *ObservableSourceMockRecorder) Has(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Has", reflect.TypeOf((*ObservableSourceMock)(nil).Has), path)
}

func (m *ObservableSourceMock) Reload() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reload")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *ObservableSourceMockRecorder) Reload() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reload", reflect.TypeOf((*ObservableSourceMock)(nil).Reload))
}

func (m *ObservableSourceMock) SetPriority(priority int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetPriority", priority)
}

func (mr *ObservableSourceMockRecorder) SetPriority(priority interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPriority", reflect.TypeOf((*ObservableSourceMock)(nil).SetPriority), priority)
}
