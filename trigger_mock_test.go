package config

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

type TriggerMock struct {
	ctrl     *gomock.Controller
	recorder *TriggerMockRecorder
}

type TriggerMockRecorder struct {
	mock *TriggerMock
}

func NewTriggerMock(ctrl *gomock.Controller) *TriggerMock {
	mock := &TriggerMock{ctrl: ctrl}
	mock.recorder = &TriggerMockRecorder{mock}
	return mock
}

func (m *TriggerMock) EXPECT() *TriggerMockRecorder {
	return m.recorder
}

func (m *TriggerMock) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *TriggerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*TriggerMock)(nil).Close))
}

func (m *TriggerMock) Delay() time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delay")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

func (mr *TriggerMockRecorder) Delay() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delay", reflect.TypeOf((*TriggerMock)(nil).Delay))
}

func (m *TriggerMock) IsClosed() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsClosed")
	ret0, _ := ret[0].(bool)
	return ret0
}

func (mr *TriggerMockRecorder) IsClosed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsClosed", reflect.TypeOf((*TriggerMock)(nil).IsClosed))
}
