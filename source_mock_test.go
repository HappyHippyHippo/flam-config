package config

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

type SourceMock struct {
	ctrl     *gomock.Controller
	recorder *SourceMockRecorder
}

type SourceMockRecorder struct {
	mock *SourceMock
}

func NewSourceMock(ctrl *gomock.Controller) *SourceMock {
	mock := &SourceMock{ctrl: ctrl}
	mock.recorder = &SourceMockRecorder{mock}
	return mock
}

func (m *SourceMock) EXPECT() *SourceMockRecorder {
	return m.recorder
}

func (m *SourceMock) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *SourceMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*SourceMock)(nil).Close))
}

func (m *SourceMock) Get(path string, def ...any) any {
	m.ctrl.T.Helper()
	varargs := append([]interface{}{path}, def...)
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0 := ret[0]
	return ret0
}

func (mr *SourceMockRecorder) Get(path interface{}, def ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{path}, def...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*SourceMock)(nil).Get), varargs...)
}

func (m *SourceMock) GetPriority() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPriority")
	ret0, _ := ret[0].(int)
	return ret0
}

func (mr *SourceMockRecorder) GetPriority() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPriority", reflect.TypeOf((*SourceMock)(nil).GetPriority))
}

func (m *SourceMock) SetPriority(priority int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetPriority", priority)
}

func (mr *SourceMockRecorder) SetPriority(priority interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPriority", reflect.TypeOf((*SourceMock)(nil).SetPriority), priority)
}
