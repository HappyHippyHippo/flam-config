package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	filesystem "github.com/happyhippyhippo/flam-filesystem"
)

// FileSystemFacade is a mock of FileSystemFacade interface.
type FileSystemFacade struct {
	ctrl     *gomock.Controller
	recorder *FileSystemFacadeRecorder
}

// FileSystemFacadeRecorder is the mock recorder for FileSystemFacade.
type FileSystemFacadeRecorder struct {
	mock *FileSystemFacade
}

// NewFileSystemFacade creates a new mock instance.
func NewFileSystemFacade(ctrl *gomock.Controller) *FileSystemFacade {
	mock := &FileSystemFacade{ctrl: ctrl}
	mock.recorder = &FileSystemFacadeRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *FileSystemFacade) EXPECT() *FileSystemFacadeRecorder {
	return m.recorder
}

// AddDisk mocks base method.
func (m *FileSystemFacade) AddDisk(id string, disk filesystem.Disk) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDisk", id, disk)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddDisk indicates an expected call of AddDisk.
func (mr *FileSystemFacadeRecorder) AddDisk(id, disk interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDisk", reflect.TypeOf((*FileSystemFacade)(nil).AddDisk), id, disk)
}

// GetDisk mocks base method.
func (m *FileSystemFacade) GetDisk(id string) (filesystem.Disk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDisk", id)
	ret0, _ := ret[0].(filesystem.Disk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDisk indicates an expected call of GetDisk.
func (mr *FileSystemFacadeRecorder) GetDisk(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDisk", reflect.TypeOf((*FileSystemFacade)(nil).GetDisk), id)
}

// HasDisk mocks base method.
func (m *FileSystemFacade) HasDisk(id string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasDisk", id)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasDisk indicates an expected call of HasDisk.
func (mr *FileSystemFacadeRecorder) HasDisk(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasDisk", reflect.TypeOf((*FileSystemFacade)(nil).HasDisk), id)
}

// ListDisks mocks base method.
func (m *FileSystemFacade) ListDisks() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListDisks")
	ret0, _ := ret[0].([]string)
	return ret0
}

// ListDisks indicates an expected call of ListDisks.
func (mr *FileSystemFacadeRecorder) ListDisks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListDisks", reflect.TypeOf((*FileSystemFacade)(nil).ListDisks))
}
