package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// Source is a mock of ConfigSource interface.
type Source struct {
	ctrl     *gomock.Controller
	recorder *SourceRecorder
}

// SourceRecorder is the mock recorder for Source.
type SourceRecorder struct {
	mock *Source
}

// NewSource creates a new mock instance.
func NewSource(ctrl *gomock.Controller) *Source {
	mock := &Source{ctrl: ctrl}
	mock.recorder = &SourceRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Source) EXPECT() *SourceRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *Source) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *SourceRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*Source)(nil).Close))
}

// Get mocks base method.
func (m *Source) Get(path string, def ...any) any {
	m.ctrl.T.Helper()
	varargs := []interface{}{path}
	for _, a := range def {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(any)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *SourceRecorder) Get(path interface{}, def ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{path}, def...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*Source)(nil).Get), varargs...)
}

// GetPriority mocks base method.
func (m *Source) GetPriority() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPriority")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetPriority indicates an expected call of GetPriority.
func (mr *SourceRecorder) GetPriority() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPriority", reflect.TypeOf((*Source)(nil).GetPriority))
}

// SetPriority mocks base method.
func (m *Source) SetPriority(priority int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetPriority", priority)
}

// SetPriority indicates an expected call of SetPriority.
func (mr *SourceRecorder) SetPriority(priority interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPriority", reflect.TypeOf((*Source)(nil).SetPriority), priority)
}
