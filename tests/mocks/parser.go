package mocks

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	flam "github.com/happyhippyhippo/flam"
)

// Parser is a mock of ConfigParser interface.
type Parser struct {
	ctrl     *gomock.Controller
	recorder *ParserRecorder
}

// ParserRecorder is the mock recorder for Parser.
type ParserRecorder struct {
	mock *Parser
}

// NewParser creates a new mock instance.
func NewParser(ctrl *gomock.Controller) *Parser {
	mock := &Parser{ctrl: ctrl}
	mock.recorder = &ParserRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Parser) EXPECT() *ParserRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *Parser) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *ParserRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*Parser)(nil).Close))
}

// Parse mocks base method.
func (m *Parser) Parse(reader io.Reader) (flam.Bag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parse", reader)
	ret0, _ := ret[0].(flam.Bag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Parse indicates an expected call of Parse.
func (mr *ParserRecorder) Parse(reader interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parse", reflect.TypeOf((*Parser)(nil).Parse), reader)
}
