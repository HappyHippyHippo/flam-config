package config

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

type RestRequesterMock struct {
	ctrl     *gomock.Controller
	recorder *RestRequesterMockRecorder
}

type RestRequesterMockRecorder struct {
	mock *RestRequesterMock
}

func NewRestRequesterMock(ctrl *gomock.Controller) *RestRequesterMock {
	mock := &RestRequesterMock{ctrl: ctrl}
	mock.recorder = &RestRequesterMockRecorder{mock}
	return mock
}

func (m *RestRequesterMock) EXPECT() *RestRequesterMockRecorder {
	return m.recorder
}

func (m *RestRequesterMock) Do(req *http.Request) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", req)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *RestRequesterMockRecorder) Do(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*RestRequesterMock)(nil).Do), req)
}
