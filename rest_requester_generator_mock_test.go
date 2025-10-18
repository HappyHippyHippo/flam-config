package config

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

type RestRequesterGeneratorMock struct {
	ctrl     *gomock.Controller
	recorder *RestRequesterGeneratorMockRecorder
}

type RestRequesterGeneratorMockRecorder struct {
	mock *RestRequesterGeneratorMock
}

func NewRestRequesterGeneratorMock(ctrl *gomock.Controller) *RestRequesterGeneratorMock {
	mock := &RestRequesterGeneratorMock{ctrl: ctrl}
	mock.recorder = &RestRequesterGeneratorMockRecorder{mock}
	return mock
}

func (m *RestRequesterGeneratorMock) EXPECT() *RestRequesterGeneratorMockRecorder {
	return m.recorder
}

func (m *RestRequesterGeneratorMock) Create() (RestRequester, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create")
	ret0, _ := ret[0].(RestRequester)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *RestRequesterGeneratorMockRecorder) Create() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*RestRequesterGeneratorMock)(nil).Create))
}
