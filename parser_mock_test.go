package config

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	flam "github.com/happyhippyhippo/flam"
)

type ParserMock struct {
	ctrl     *gomock.Controller
	recorder *ParserMockRecorder
}

type ParserMockRecorder struct {
	mock *ParserMock
}

func NewParserMock(ctrl *gomock.Controller) *ParserMock {
	mock := &ParserMock{ctrl: ctrl}
	mock.recorder = &ParserMockRecorder{mock}
	return mock
}

func (m *ParserMock) EXPECT() *ParserMockRecorder {
	return m.recorder
}

func (m *ParserMock) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *ParserMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*ParserMock)(nil).Close))
}

func (m *ParserMock) Parse(reader io.Reader) (flam.Bag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parse", reader)
	ret0, _ := ret[0].(flam.Bag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *ParserMockRecorder) Parse(reader interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parse", reflect.TypeOf((*ParserMock)(nil).Parse), reader)
}
