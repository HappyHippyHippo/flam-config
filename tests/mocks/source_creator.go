package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	flam "github.com/happyhippyhippo/flam"
	config "github.com/happyhippyhippo/flam-config"
)

// SourceCreator is a mock of ConfigSourceCreator interface.
type SourceCreator struct {
	ctrl     *gomock.Controller
	recorder *SourceCreatorRecorder
}

// SourceCreatorRecorder is the mock recorder for SourceCreator.
type SourceCreatorRecorder struct {
	mock *SourceCreator
}

// NewSourceCreator creates a new mock instance.
func NewSourceCreator(ctrl *gomock.Controller) *SourceCreator {
	mock := &SourceCreator{ctrl: ctrl}
	mock.recorder = &SourceCreatorRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *SourceCreator) EXPECT() *SourceCreatorRecorder {
	return m.recorder
}

// Accept mocks base method.
func (m *SourceCreator) Accept(config flam.Bag) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accept", config)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Accept indicates an expected call of Accept.
func (mr *SourceCreatorRecorder) Accept(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*SourceCreator)(nil).Accept), config)
}

// Create mocks base method.
func (m *SourceCreator) Create(cfg flam.Bag) (config.Source, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", cfg)
	ret0, _ := ret[0].(config.Source)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *SourceCreatorRecorder) Create(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*SourceCreator)(nil).Create), config)
}
