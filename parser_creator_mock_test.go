package config

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	flam "github.com/happyhippyhippo/flam"
)

type ParserCreatorMock struct {
	ctrl     *gomock.Controller
	recorder *ParserCreatorMockRecorder
}

type ParserCreatorMockRecorder struct {
	mock *ParserCreatorMock
}

func NewParserCreatorMock(ctrl *gomock.Controller) *ParserCreatorMock {
	mock := &ParserCreatorMock{ctrl: ctrl}
	mock.recorder = &ParserCreatorMockRecorder{mock}
	return mock
}

func (m *ParserCreatorMock) EXPECT() *ParserCreatorMockRecorder {
	return m.recorder
}

func (m *ParserCreatorMock) Accept(config flam.Bag) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accept", config)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (mr *ParserCreatorMockRecorder) Accept(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*ParserCreatorMock)(nil).Accept), config)
}

func (m *ParserCreatorMock) Create(cfg flam.Bag) (Parser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", cfg)
	ret0, _ := ret[0].(Parser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *ParserCreatorMockRecorder) Create(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*ParserCreatorMock)(nil).Create), config)
}
