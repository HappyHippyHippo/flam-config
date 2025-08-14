package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	flam "github.com/happyhippyhippo/flam"
	config "github.com/happyhippyhippo/flam-config"
)

// ParserCreator is a mock of ConfigParserCreator interface.
type ParserCreator struct {
	ctrl     *gomock.Controller
	recorder *ParserCreatorRecorder
}

// ParserCreatorRecorder is the mock recorder for ParserCreator.
type ParserCreatorRecorder struct {
	mock *ParserCreator
}

// NewParserCreator creates a new mock instance.
func NewParserCreator(ctrl *gomock.Controller) *ParserCreator {
	mock := &ParserCreator{ctrl: ctrl}
	mock.recorder = &ParserCreatorRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *ParserCreator) EXPECT() *ParserCreatorRecorder {
	return m.recorder
}

// Accept mocks base method.
func (m *ParserCreator) Accept(config flam.Bag) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accept", config)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Accept indicates an expected call of Accept.
func (mr *ParserCreatorRecorder) Accept(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*ParserCreator)(nil).Accept), config)
}

// Create mocks base method.
func (m *ParserCreator) Create(cfg flam.Bag) (config.Parser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", cfg)
	ret0, _ := ret[0].(config.Parser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *ParserCreatorRecorder) Create(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*ParserCreator)(nil).Create), config)
}
