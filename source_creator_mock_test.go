package config

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	flam "github.com/happyhippyhippo/flam"
)

type SourceCreatorMock struct {
	ctrl     *gomock.Controller
	recorder *SourceCreatorMockRecorder
}

type SourceCreatorMockRecorder struct {
	mock *SourceCreatorMock
}

func NewSourceCreatorMock(ctrl *gomock.Controller) *SourceCreatorMock {
	mock := &SourceCreatorMock{ctrl: ctrl}
	mock.recorder = &SourceCreatorMockRecorder{mock}
	return mock
}

func (m *SourceCreatorMock) EXPECT() *SourceCreatorMockRecorder {
	return m.recorder
}

func (m *SourceCreatorMock) Accept(config flam.Bag) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accept", config)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (mr *SourceCreatorMockRecorder) Accept(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*SourceCreatorMock)(nil).Accept), config)
}

func (m *SourceCreatorMock) Create(cfg flam.Bag) (Source, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", cfg)
	ret0, _ := ret[0].(Source)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *SourceCreatorMockRecorder) Create(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*SourceCreatorMock)(nil).Create), config)
}
