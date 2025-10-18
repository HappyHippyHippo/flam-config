package config

import (
	os "os"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	afero "github.com/spf13/afero"
)

type DiskMock struct {
	ctrl     *gomock.Controller
	recorder *DiskMockRecorder
}

type DiskMockRecorder struct {
	mock *DiskMock
}

func NewDiskMock(ctrl *gomock.Controller) *DiskMock {
	mock := &DiskMock{ctrl: ctrl}
	mock.recorder = &DiskMockRecorder{mock}
	return mock
}

func (m *DiskMock) EXPECT() *DiskMockRecorder {
	return m.recorder
}

func (m *DiskMock) Chmod(name string, mode os.FileMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Chmod", name, mode)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *DiskMockRecorder) Chmod(name, mode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Chmod", reflect.TypeOf((*DiskMock)(nil).Chmod), name, mode)
}

func (m *DiskMock) Chown(name string, uid, gid int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Chown", name, uid, gid)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *DiskMockRecorder) Chown(name, uid, gid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Chown", reflect.TypeOf((*DiskMock)(nil).Chown), name, uid, gid)
}

func (m *DiskMock) Chtimes(name string, atime, mtime time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Chtimes", name, atime, mtime)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *DiskMockRecorder) Chtimes(name, atime, mtime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Chtimes", reflect.TypeOf((*DiskMock)(nil).Chtimes), name, atime, mtime)
}

func (m *DiskMock) Create(name string) (afero.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", name)
	ret0, _ := ret[0].(afero.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *DiskMockRecorder) Create(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*DiskMock)(nil).Create), name)
}

func (m *DiskMock) Mkdir(name string, perm os.FileMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Mkdir", name, perm)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *DiskMockRecorder) Mkdir(name, perm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Mkdir", reflect.TypeOf((*DiskMock)(nil).Mkdir), name, perm)
}

func (m *DiskMock) MkdirAll(path string, perm os.FileMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MkdirAll", path, perm)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *DiskMockRecorder) MkdirAll(path, perm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MkdirAll", reflect.TypeOf((*DiskMock)(nil).MkdirAll), path, perm)
}

func (m *DiskMock) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

func (mr *DiskMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*DiskMock)(nil).Name))
}

func (m *DiskMock) Open(name string) (afero.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open", name)
	ret0, _ := ret[0].(afero.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *DiskMockRecorder) Open(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*DiskMock)(nil).Open), name)
}

func (m *DiskMock) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenFile", name, flag, perm)
	ret0, _ := ret[0].(afero.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *DiskMockRecorder) OpenFile(name, flag, perm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenFile", reflect.TypeOf((*DiskMock)(nil).OpenFile), name, flag, perm)
}

func (m *DiskMock) Remove(name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", name)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *DiskMockRecorder) Remove(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*DiskMock)(nil).Remove), name)
}

func (m *DiskMock) RemoveAll(path string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAll", path)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *DiskMockRecorder) RemoveAll(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAll", reflect.TypeOf((*DiskMock)(nil).RemoveAll), path)
}

func (m *DiskMock) Rename(oldname, newname string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rename", oldname, newname)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *DiskMockRecorder) Rename(oldname, newname interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rename", reflect.TypeOf((*DiskMock)(nil).Rename), oldname, newname)
}

func (m *DiskMock) Stat(name string) (os.FileInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stat", name)
	ret0, _ := ret[0].(os.FileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *DiskMockRecorder) Stat(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stat", reflect.TypeOf((*DiskMock)(nil).Stat), name)
}
