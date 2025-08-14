package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	config "github.com/happyhippyhippo/flam-config"
)

// RestRequesterGenerator is a mock of ConfigRestRequesterGenerator interface.
type RestRequesterGenerator struct {
	ctrl     *gomock.Controller
	recorder *RestRequesterGeneratorRecorder
}

// RestRequesterGeneratorRecorder is the mock recorder for RestRequesterGenerator.
type RestRequesterGeneratorRecorder struct {
	mock *RestRequesterGenerator
}

// NewRestRequesterGenerator creates a new mock instance.
func NewRestRequesterGenerator(ctrl *gomock.Controller) *RestRequesterGenerator {
	mock := &RestRequesterGenerator{ctrl: ctrl}
	mock.recorder = &RestRequesterGeneratorRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *RestRequesterGenerator) EXPECT() *RestRequesterGeneratorRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *RestRequesterGenerator) Create() (config.RestRequester, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create")
	ret0, _ := ret[0].(config.RestRequester)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *RestRequesterGeneratorRecorder) Create() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*RestRequesterGenerator)(nil).Create))
}
