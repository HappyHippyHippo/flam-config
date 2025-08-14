package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// Trigger is a mock of Trigger interface.
type Trigger struct {
	ctrl     *gomock.Controller
	recorder *TriggerRecorder
}

// TriggerRecorder is the mock recorder for Trigger.
type TriggerRecorder struct {
	mock *Trigger
}

// NewTrigger creates a new mock instance.
func NewTrigger(ctrl *gomock.Controller) *Trigger {
	mock := &Trigger{ctrl: ctrl}
	mock.recorder = &TriggerRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Trigger) EXPECT() *TriggerRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *Trigger) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *TriggerRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*Trigger)(nil).Close))
}

// Delay mocks base method.
func (m *Trigger) Delay() time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delay")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// Delay indicates an expected call of Delay.
func (mr *TriggerRecorder) Delay() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delay", reflect.TypeOf((*Trigger)(nil).Delay))
}

// IsClosed mocks base method.
func (m *Trigger) IsClosed() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsClosed")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsClosed indicates an expected call of IsClosed.
func (mr *TriggerRecorder) IsClosed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsClosed", reflect.TypeOf((*Trigger)(nil).IsClosed))
}
