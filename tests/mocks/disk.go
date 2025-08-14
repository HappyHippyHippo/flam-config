package mocks

import (
	os "os"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	afero "github.com/spf13/afero"
)

// Disk is a mock of FileSystemDisk interface.
type Disk struct {
	ctrl     *gomock.Controller
	recorder *DiskRecorder
}

// DiskRecorder is the mock recorder for Disk.
type DiskRecorder struct {
	mock *Disk
}

// NewDisk creates a new mock instance.
func NewDisk(ctrl *gomock.Controller) *Disk {
	mock := &Disk{ctrl: ctrl}
	mock.recorder = &DiskRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Disk) EXPECT() *DiskRecorder {
	return m.recorder
}

// Chmod mocks base method.
func (m *Disk) Chmod(name string, mode os.FileMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Chmod", name, mode)
	ret0, _ := ret[0].(error)
	return ret0
}

// Chmod indicates an expected call of Chmod.
func (mr *DiskRecorder) Chmod(name, mode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Chmod", reflect.TypeOf((*Disk)(nil).Chmod), name, mode)
}

// Chown mocks base method.
func (m *Disk) Chown(name string, uid, gid int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Chown", name, uid, gid)
	ret0, _ := ret[0].(error)
	return ret0
}

// Chown indicates an expected call of Chown.
func (mr *DiskRecorder) Chown(name, uid, gid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Chown", reflect.TypeOf((*Disk)(nil).Chown), name, uid, gid)
}

// Chtimes mocks base method.
func (m *Disk) Chtimes(name string, atime, mtime time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Chtimes", name, atime, mtime)
	ret0, _ := ret[0].(error)
	return ret0
}

// Chtimes indicates an expected call of Chtimes.
func (mr *DiskRecorder) Chtimes(name, atime, mtime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Chtimes", reflect.TypeOf((*Disk)(nil).Chtimes), name, atime, mtime)
}

// Create mocks base method.
func (m *Disk) Create(name string) (afero.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", name)
	ret0, _ := ret[0].(afero.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *DiskRecorder) Create(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*Disk)(nil).Create), name)
}

// Mkdir mocks base method.
func (m *Disk) Mkdir(name string, perm os.FileMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Mkdir", name, perm)
	ret0, _ := ret[0].(error)
	return ret0
}

// Mkdir indicates an expected call of Mkdir.
func (mr *DiskRecorder) Mkdir(name, perm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Mkdir", reflect.TypeOf((*Disk)(nil).Mkdir), name, perm)
}

// MkdirAll mocks base method.
func (m *Disk) MkdirAll(path string, perm os.FileMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MkdirAll", path, perm)
	ret0, _ := ret[0].(error)
	return ret0
}

// MkdirAll indicates an expected call of MkdirAll.
func (mr *DiskRecorder) MkdirAll(path, perm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MkdirAll", reflect.TypeOf((*Disk)(nil).MkdirAll), path, perm)
}

// Name mocks base method.
func (m *Disk) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *DiskRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*Disk)(nil).Name))
}

// Open mocks base method.
func (m *Disk) Open(name string) (afero.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open", name)
	ret0, _ := ret[0].(afero.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Open indicates an expected call of Open.
func (mr *DiskRecorder) Open(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*Disk)(nil).Open), name)
}

// OpenFile mocks base method.
func (m *Disk) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenFile", name, flag, perm)
	ret0, _ := ret[0].(afero.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenFile indicates an expected call of OpenFile.
func (mr *DiskRecorder) OpenFile(name, flag, perm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenFile", reflect.TypeOf((*Disk)(nil).OpenFile), name, flag, perm)
}

// Remove mocks base method.
func (m *Disk) Remove(name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", name)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *DiskRecorder) Remove(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*Disk)(nil).Remove), name)
}

// RemoveAll mocks base method.
func (m *Disk) RemoveAll(path string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAll", path)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveAll indicates an expected call of RemoveAll.
func (mr *DiskRecorder) RemoveAll(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAll", reflect.TypeOf((*Disk)(nil).RemoveAll), path)
}

// Rename mocks base method.
func (m *Disk) Rename(oldname, newname string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rename", oldname, newname)
	ret0, _ := ret[0].(error)
	return ret0
}

// Rename indicates an expected call of Rename.
func (mr *DiskRecorder) Rename(oldname, newname interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rename", reflect.TypeOf((*Disk)(nil).Rename), oldname, newname)
}

// Stat mocks base method.
func (m *Disk) Stat(name string) (os.FileInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stat", name)
	ret0, _ := ret[0].(os.FileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stat indicates an expected call of Stat.
func (mr *DiskRecorder) Stat(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stat", reflect.TypeOf((*Disk)(nil).Stat), name)
}
