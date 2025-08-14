package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// ObservableSource is a mock of ConfigObservableSource interface.
type ObservableSource struct {
	ctrl     *gomock.Controller
	recorder *ObservableSourceRecorder
}

// ObservableSourceRecorder is the mock recorder for ObservableSource.
type ObservableSourceRecorder struct {
	mock *ObservableSource
}

// NewObservableSource creates a new mock instance.
func NewObservableSource(ctrl *gomock.Controller) *ObservableSource {
	mock := &ObservableSource{ctrl: ctrl}
	mock.recorder = &ObservableSourceRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *ObservableSource) EXPECT() *ObservableSourceRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *ObservableSource) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *ObservableSourceRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*ObservableSource)(nil).Close))
}

// Get mocks base method.
func (m *ObservableSource) Get(path string, def ...any) any {
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
func (mr *ObservableSourceRecorder) Get(path interface{}, def ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{path}, def...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*ObservableSource)(nil).Get), varargs...)
}

// GetPriority mocks base method.
func (m *ObservableSource) GetPriority() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPriority")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetPriority indicates an expected call of GetPriority.
func (mr *ObservableSourceRecorder) GetPriority() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPriority", reflect.TypeOf((*ObservableSource)(nil).GetPriority))
}

// Has mocks base method.
func (m *ObservableSource) Has(path string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Has", path)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Has indicates an expected call of Has.
func (mr *ObservableSourceRecorder) Has(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Has", reflect.TypeOf((*ObservableSource)(nil).Has), path)
}

// Reload mocks base method.
func (m *ObservableSource) Reload() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reload")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Reload indicates an expected call of Reload.
func (mr *ObservableSourceRecorder) Reload() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reload", reflect.TypeOf((*ObservableSource)(nil).Reload))
}

// SetPriority mocks base method.
func (m *ObservableSource) SetPriority(priority int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetPriority", priority)
}

// SetPriority indicates an expected call of SetPriority.
func (mr *ObservableSourceRecorder) SetPriority(priority interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPriority", reflect.TypeOf((*ObservableSource)(nil).SetPriority), priority)
}
