package config

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

type ReadCloserMock struct {
	ctrl     *gomock.Controller
	recorder *ReadCloserMockRecorder
}

type ReadCloserMockRecorder struct {
	mock *ReadCloserMock
}

func NewReadCloserMock(ctrl *gomock.Controller) *ReadCloserMock {
	mock := &ReadCloserMock{ctrl: ctrl}
	mock.recorder = &ReadCloserMockRecorder{mock}
	return mock
}

func (m *ReadCloserMock) EXPECT() *ReadCloserMockRecorder {
	return m.recorder
}

func (m *ReadCloserMock) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *ReadCloserMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*ReadCloserMock)(nil).Close))
}

func (m *ReadCloserMock) Read(arg0 []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *ReadCloserMockRecorder) Read(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*ReadCloserMock)(nil).Read), arg0)
}
