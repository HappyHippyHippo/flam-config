package config

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	flam "github.com/happyhippyhippo/flam"
)

type FactoryConfigMock struct {
	ctrl     *gomock.Controller
	recorder *FactoryConfigMockRecorder
}

type FactoryConfigMockRecorder struct {
	mock *FactoryConfigMock
}

func NewFactoryConfigMock(ctrl *gomock.Controller) *FactoryConfigMock {
	mock := &FactoryConfigMock{ctrl: ctrl}
	mock.recorder = &FactoryConfigMockRecorder{mock}
	return mock
}

func (m *FactoryConfigMock) EXPECT() *FactoryConfigMockRecorder {
	return m.recorder
}

func (m *FactoryConfigMock) Get(path string, def ...any) flam.Bag {
	m.ctrl.T.Helper()
	varargs := append([]any{path}, def...)
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(flam.Bag)
	return ret0
}

func (mr *FactoryConfigMockRecorder) Get(path any, def ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{path}, def...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*FactoryConfigMock)(nil).Get), varargs...)
}
