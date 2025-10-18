package config

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	flam "github.com/happyhippyhippo/flam"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
)

type DiskCreatorMock struct {
	ctrl     *gomock.Controller
	recorder *DiskCreatorMockRecorder
}

type DiskCreatorMockRecorder struct {
	mock *DiskCreatorMock
}

func NewDiskCreatorMock(ctrl *gomock.Controller) *DiskCreatorMock {
	mock := &DiskCreatorMock{ctrl: ctrl}
	mock.recorder = &DiskCreatorMockRecorder{mock}
	return mock
}

func (m *DiskCreatorMock) EXPECT() *DiskCreatorMockRecorder {
	return m.recorder
}

func (m *DiskCreatorMock) Accept(config flam.Bag) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accept", config)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (mr *DiskCreatorMockRecorder) Accept(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*DiskCreatorMock)(nil).Accept), config)
}

func (m *DiskCreatorMock) Create(config flam.Bag) (filesystem.Disk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", config)
	ret0, _ := ret[0].(filesystem.Disk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *DiskCreatorMockRecorder) Create(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*DiskCreatorMock)(nil).Create), config)
}
