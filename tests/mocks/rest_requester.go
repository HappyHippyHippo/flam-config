package mocks

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// RestRequester is a mock of RestRequester interface.
type RestRequester struct {
	ctrl     *gomock.Controller
	recorder *RestRequesterRecorder
}

// RestRequesterRecorder is the mock recorder for RestRequester.
type RestRequesterRecorder struct {
	mock *RestRequester
}

// NewRestRequester creates a new mock instance.
func NewRestRequester(ctrl *gomock.Controller) *RestRequester {
	mock := &RestRequester{ctrl: ctrl}
	mock.recorder = &RestRequesterRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *RestRequester) EXPECT() *RestRequesterRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *RestRequester) Do(req *http.Request) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", req)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *RestRequesterRecorder) Do(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*RestRequester)(nil).Do), req)
}
