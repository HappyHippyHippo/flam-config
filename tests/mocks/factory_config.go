package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	flam "github.com/happyhippyhippo/flam"
)

// FactoryConfig is a mock of FactoryConfig interface.
type FactoryConfig struct {
	ctrl     *gomock.Controller
	recorder *FactoryConfigRecorder
}

// FactoryConfigRecorder is the mock recorder for FactoryConfig.
type FactoryConfigRecorder struct {
	mock *FactoryConfig
}

// NewFactoryConfig creates a new mock instance.
func NewFactoryConfig(ctrl *gomock.Controller) *FactoryConfig {
	mock := &FactoryConfig{ctrl: ctrl}
	mock.recorder = &FactoryConfigRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *FactoryConfig) EXPECT() *FactoryConfigRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *FactoryConfig) Get(path string, def ...any) flam.Bag {
	m.ctrl.T.Helper()
	varargs := []any{path}
	for _, a := range def {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(flam.Bag)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *FactoryConfigRecorder) Get(path any, def ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{path}, def...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*FactoryConfig)(nil).Get), varargs...)
}
