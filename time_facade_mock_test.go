package config

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"

	flamTime "github.com/happyhippyhippo/flam-time"
)

type TimeFacadeMock struct {
	ctrl     *gomock.Controller
	recorder *TimeFacadeMockRecorder
}

type TimeFacadeMockRecorder struct {
	mock *TimeFacadeMock
}

func NewTimeFacadeMock(ctrl *gomock.Controller) *TimeFacadeMock {
	mock := &TimeFacadeMock{ctrl: ctrl}
	mock.recorder = &TimeFacadeMockRecorder{mock}
	return mock
}

func (m *TimeFacadeMock) EXPECT() *TimeFacadeMockRecorder {
	return m.recorder
}

func (m *TimeFacadeMock) Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Date", year, month, day, hour, min, sec, nsec, loc)
	ret0, _ := ret[0].(time.Time)
	return ret0
}

func (mr *TimeFacadeMockRecorder) Date(year, month, day, hour, min, sec, nsec, loc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Date", reflect.TypeOf((*TimeFacadeMock)(nil).Date), year, month, day, hour, min, sec, nsec, loc)
}

func (m *TimeFacadeMock) FixedZone(name string, offset int) *time.Location {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FixedZone", name, offset)
	ret0, _ := ret[0].(*time.Location)
	return ret0
}

func (mr *TimeFacadeMockRecorder) FixedZone(name, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FixedZone", reflect.TypeOf((*TimeFacadeMock)(nil).FixedZone), name, offset)
}

func (m *TimeFacadeMock) LoadLocation(name string) (*time.Location, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadLocation", name)
	ret0, _ := ret[0].(*time.Location)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *TimeFacadeMockRecorder) LoadLocation(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadLocation", reflect.TypeOf((*TimeFacadeMock)(nil).LoadLocation), name)
}

func (m *TimeFacadeMock) NewPulseTrigger(delay time.Duration, callback flamTime.Callback) (flamTime.Trigger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewPulseTrigger", delay, callback)
	ret0, _ := ret[0].(flamTime.Trigger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *TimeFacadeMockRecorder) NewPulseTrigger(delay, callback interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewPulseTrigger", reflect.TypeOf((*TimeFacadeMock)(nil).NewPulseTrigger), delay, callback)
}

func (m *TimeFacadeMock) NewRecurringTrigger(delay time.Duration, callback flamTime.Callback) (flamTime.Trigger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRecurringTrigger", delay, callback)
	ret0, _ := ret[0].(flamTime.Trigger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *TimeFacadeMockRecorder) NewRecurringTrigger(delay, callback interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRecurringTrigger", reflect.TypeOf((*TimeFacadeMock)(nil).NewRecurringTrigger), delay, callback)
}

func (m *TimeFacadeMock) Now() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Now")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

func (mr *TimeFacadeMockRecorder) Now() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Now", reflect.TypeOf((*TimeFacadeMock)(nil).Now))
}

func (m *TimeFacadeMock) Parse(layout, value string) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parse", layout, value)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *TimeFacadeMockRecorder) Parse(layout, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parse", reflect.TypeOf((*TimeFacadeMock)(nil).Parse), layout, value)
}

func (m *TimeFacadeMock) ParseDuration(s string) (time.Duration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseDuration", s)
	ret0, _ := ret[0].(time.Duration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *TimeFacadeMockRecorder) ParseDuration(s interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseDuration", reflect.TypeOf((*TimeFacadeMock)(nil).ParseDuration), s)
}

func (m *TimeFacadeMock) ParseInLocation(layout, value string, loc *time.Location) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseInLocation", layout, value, loc)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *TimeFacadeMockRecorder) ParseInLocation(layout, value, loc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseInLocation", reflect.TypeOf((*TimeFacadeMock)(nil).ParseInLocation), layout, value, loc)
}

func (m *TimeFacadeMock) Since(t time.Time) time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Since", t)
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

func (mr *TimeFacadeMockRecorder) Since(t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Since", reflect.TypeOf((*TimeFacadeMock)(nil).Since), t)
}

func (m *TimeFacadeMock) Unix(sec, nsec int64) time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unix", sec, nsec)
	ret0, _ := ret[0].(time.Time)
	return ret0
}

func (mr *TimeFacadeMockRecorder) Unix(sec, nsec interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unix", reflect.TypeOf((*TimeFacadeMock)(nil).Unix), sec, nsec)
}

func (m *TimeFacadeMock) UnixMicro(usec int64) time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnixMicro", usec)
	ret0, _ := ret[0].(time.Time)
	return ret0
}

func (mr *TimeFacadeMockRecorder) UnixMicro(usec interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnixMicro", reflect.TypeOf((*TimeFacadeMock)(nil).UnixMicro), usec)
}

func (m *TimeFacadeMock) UnixMilli(msec int64) time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnixMilli", msec)
	ret0, _ := ret[0].(time.Time)
	return ret0
}

func (mr *TimeFacadeMockRecorder) UnixMilli(msec interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnixMilli", reflect.TypeOf((*TimeFacadeMock)(nil).UnixMilli), msec)
}

func (m *TimeFacadeMock) Until(t time.Time) time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Until", t)
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

func (mr *TimeFacadeMockRecorder) Until(t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Until", reflect.TypeOf((*TimeFacadeMock)(nil).Until), t)
}
