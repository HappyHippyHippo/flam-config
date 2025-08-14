package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"

	flamTime "github.com/happyhippyhippo/flam-time"
)

// TimeFacade is a mock of TimeFacade interface.
type TimeFacade struct {
	ctrl     *gomock.Controller
	recorder *TimeFacadeRecorder
}

// TimeFacadeRecorder is the mock recorder for TimeFacade.
type TimeFacadeRecorder struct {
	mock *TimeFacade
}

// NewTimeFacade creates a new mock instance.
func NewTimeFacade(ctrl *gomock.Controller) *TimeFacade {
	mock := &TimeFacade{ctrl: ctrl}
	mock.recorder = &TimeFacadeRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *TimeFacade) EXPECT() *TimeFacadeRecorder {
	return m.recorder
}

// Date mocks base method.
func (m *TimeFacade) Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Date", year, month, day, hour, min, sec, nsec, loc)
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Date indicates an expected call of Date.
func (mr *TimeFacadeRecorder) Date(year, month, day, hour, min, sec, nsec, loc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Date", reflect.TypeOf((*TimeFacade)(nil).Date), year, month, day, hour, min, sec, nsec, loc)
}

// FixedZone mocks base method.
func (m *TimeFacade) FixedZone(name string, offset int) *time.Location {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FixedZone", name, offset)
	ret0, _ := ret[0].(*time.Location)
	return ret0
}

// FixedZone indicates an expected call of FixedZone.
func (mr *TimeFacadeRecorder) FixedZone(name, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FixedZone", reflect.TypeOf((*TimeFacade)(nil).FixedZone), name, offset)
}

// LoadLocation mocks base method.
func (m *TimeFacade) LoadLocation(name string) (*time.Location, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadLocation", name)
	ret0, _ := ret[0].(*time.Location)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadLocation indicates an expected call of LoadLocation.
func (mr *TimeFacadeRecorder) LoadLocation(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadLocation", reflect.TypeOf((*TimeFacade)(nil).LoadLocation), name)
}

// NewPulseTrigger mocks base method.
func (m *TimeFacade) NewPulseTrigger(delay time.Duration, callback flamTime.Callback) (flamTime.Trigger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewPulseTrigger", delay, callback)
	ret0, _ := ret[0].(flamTime.Trigger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewPulseTrigger indicates an expected call of NewPulseTrigger.
func (mr *TimeFacadeRecorder) NewPulseTrigger(delay, callback interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewPulseTrigger", reflect.TypeOf((*TimeFacade)(nil).NewPulseTrigger), delay, callback)
}

// NewRecurringTrigger mocks base method.
func (m *TimeFacade) NewRecurringTrigger(delay time.Duration, callback flamTime.Callback) (flamTime.Trigger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRecurringTrigger", delay, callback)
	ret0, _ := ret[0].(flamTime.Trigger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewRecurringTrigger indicates an expected call of NewRecurringTrigger.
func (mr *TimeFacadeRecorder) NewRecurringTrigger(delay, callback interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRecurringTrigger", reflect.TypeOf((*TimeFacade)(nil).NewRecurringTrigger), delay, callback)
}

// Now mocks base method.
func (m *TimeFacade) Now() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Now")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Now indicates an expected call of Now.
func (mr *TimeFacadeRecorder) Now() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Now", reflect.TypeOf((*TimeFacade)(nil).Now))
}

// Parse mocks base method.
func (m *TimeFacade) Parse(layout, value string) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parse", layout, value)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Parse indicates an expected call of Parse.
func (mr *TimeFacadeRecorder) Parse(layout, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parse", reflect.TypeOf((*TimeFacade)(nil).Parse), layout, value)
}

// ParseDuration mocks base method.
func (m *TimeFacade) ParseDuration(s string) (time.Duration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseDuration", s)
	ret0, _ := ret[0].(time.Duration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseDuration indicates an expected call of ParseDuration.
func (mr *TimeFacadeRecorder) ParseDuration(s interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseDuration", reflect.TypeOf((*TimeFacade)(nil).ParseDuration), s)
}

// ParseInLocation mocks base method.
func (m *TimeFacade) ParseInLocation(layout, value string, loc *time.Location) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseInLocation", layout, value, loc)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseInLocation indicates an expected call of ParseInLocation.
func (mr *TimeFacadeRecorder) ParseInLocation(layout, value, loc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseInLocation", reflect.TypeOf((*TimeFacade)(nil).ParseInLocation), layout, value, loc)
}

// Since mocks base method.
func (m *TimeFacade) Since(t time.Time) time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Since", t)
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// Since indicates an expected call of Since.
func (mr *TimeFacadeRecorder) Since(t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Since", reflect.TypeOf((*TimeFacade)(nil).Since), t)
}

// Unix mocks base method.
func (m *TimeFacade) Unix(sec, nsec int64) time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unix", sec, nsec)
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Unix indicates an expected call of Unix.
func (mr *TimeFacadeRecorder) Unix(sec, nsec interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unix", reflect.TypeOf((*TimeFacade)(nil).Unix), sec, nsec)
}

// UnixMicro mocks base method.
func (m *TimeFacade) UnixMicro(usec int64) time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnixMicro", usec)
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// UnixMicro indicates an expected call of UnixMicro.
func (mr *TimeFacadeRecorder) UnixMicro(usec interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnixMicro", reflect.TypeOf((*TimeFacade)(nil).UnixMicro), usec)
}

// UnixMilli mocks base method.
func (m *TimeFacade) UnixMilli(msec int64) time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnixMilli", msec)
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// UnixMilli indicates an expected call of UnixMilli.
func (mr *TimeFacadeRecorder) UnixMilli(msec interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnixMilli", reflect.TypeOf((*TimeFacade)(nil).UnixMilli), msec)
}

// Until mocks base method.
func (m *TimeFacade) Until(t time.Time) time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Until", t)
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// Until indicates an expected call of Until.
func (mr *TimeFacadeRecorder) Until(t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Until", reflect.TypeOf((*TimeFacade)(nil).Until), t)
}
