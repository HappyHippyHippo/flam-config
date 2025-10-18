package config

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	filesystem "github.com/happyhippyhippo/flam-filesystem"
)

type FileSystemFacadeMock struct {
	ctrl     *gomock.Controller
	recorder *FileSystemFacadeMockRecorder
}

type FileSystemFacadeMockRecorder struct {
	mock *FileSystemFacadeMock
}

func NewFileSystemFacadeMock(ctrl *gomock.Controller) *FileSystemFacadeMock {
	mock := &FileSystemFacadeMock{ctrl: ctrl}
	mock.recorder = &FileSystemFacadeMockRecorder{mock}
	return mock
}

func (m *FileSystemFacadeMock) EXPECT() *FileSystemFacadeMockRecorder {
	return m.recorder
}

func (m *FileSystemFacadeMock) AddDisk(id string, disk filesystem.Disk) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDisk", id, disk)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *FileSystemFacadeMockRecorder) AddDisk(id, disk interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDisk", reflect.TypeOf((*FileSystemFacadeMock)(nil).AddDisk), id, disk)
}

func (m *FileSystemFacadeMock) GetDisk(id string) (filesystem.Disk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDisk", id)
	ret0, _ := ret[0].(filesystem.Disk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *FileSystemFacadeMockRecorder) GetDisk(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDisk", reflect.TypeOf((*FileSystemFacadeMock)(nil).GetDisk), id)
}

func (m *FileSystemFacadeMock) HasDisk(id string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasDisk", id)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (mr *FileSystemFacadeMockRecorder) HasDisk(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasDisk", reflect.TypeOf((*FileSystemFacadeMock)(nil).HasDisk), id)
}

func (m *FileSystemFacadeMock) ListDisks() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListDisks")
	ret0, _ := ret[0].([]string)
	return ret0
}

func (mr *FileSystemFacadeMockRecorder) ListDisks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListDisks", reflect.TypeOf((*FileSystemFacadeMock)(nil).ListDisks))
}
