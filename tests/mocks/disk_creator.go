package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	flam "github.com/happyhippyhippo/flam"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
)

// DiskCreator is a mock of FileSystemDiskCreator interface.
type DiskCreator struct {
	ctrl     *gomock.Controller
	recorder *DiskCreatorRecorder
}

// DiskCreatorRecorder is the mock recorder for DiskCreator.
type DiskCreatorRecorder struct {
	mock *DiskCreator
}

// NewDiskCreator creates a new mock instance.
func NewDiskCreator(ctrl *gomock.Controller) *DiskCreator {
	mock := &DiskCreator{ctrl: ctrl}
	mock.recorder = &DiskCreatorRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *DiskCreator) EXPECT() *DiskCreatorRecorder {
	return m.recorder
}

// Accept mocks base method.
func (m *DiskCreator) Accept(config flam.Bag) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accept", config)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Accept indicates an expected call of Accept.
func (mr *DiskCreatorRecorder) Accept(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*DiskCreator)(nil).Accept), config)
}

// Create mocks base method.
func (m *DiskCreator) Create(config flam.Bag) (filesystem.Disk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", config)
	ret0, _ := ret[0].(filesystem.Disk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *DiskCreatorRecorder) Create(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*DiskCreator)(nil).Create), config)
}
