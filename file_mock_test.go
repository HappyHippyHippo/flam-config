package config

import (
	os "os"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

type FileMock struct {
	ctrl     *gomock.Controller
	recorder *FileMockRecorder
}

type FileMockRecorder struct {
	mock *FileMock
}

func NewFileMock(ctrl *gomock.Controller) *FileMock {
	mock := &FileMock{ctrl: ctrl}
	mock.recorder = &FileMockRecorder{mock}
	return mock
}

func (m *FileMock) EXPECT() *FileMockRecorder {
	return m.recorder
}

func (m *FileMock) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *FileMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*FileMock)(nil).Close))
}

func (m *FileMock) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

func (mr *FileMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*FileMock)(nil).Name))
}

func (m *FileMock) Read(p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *FileMockRecorder) Read(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*FileMock)(nil).Read), p)
}

func (m *FileMock) ReadAt(p []byte, off int64) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadAt", p, off)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *FileMockRecorder) ReadAt(p, off interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAt", reflect.TypeOf((*FileMock)(nil).ReadAt), p, off)
}

func (m *FileMock) Readdir(count int) ([]os.FileInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Readdir", count)
	ret0, _ := ret[0].([]os.FileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *FileMockRecorder) Readdir(count interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Readdir", reflect.TypeOf((*FileMock)(nil).Readdir), count)
}

func (m *FileMock) Readdirnames(n int) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Readdirnames", n)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *FileMockRecorder) Readdirnames(n interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Readdirnames", reflect.TypeOf((*FileMock)(nil).Readdirnames), n)
}

func (m *FileMock) Seek(offset int64, whence int) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Seek", offset, whence)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *FileMockRecorder) Seek(offset, whence interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Seek", reflect.TypeOf((*FileMock)(nil).Seek), offset, whence)
}

func (m *FileMock) Stat() (os.FileInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stat")
	ret0, _ := ret[0].(os.FileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *FileMockRecorder) Stat() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stat", reflect.TypeOf((*FileMock)(nil).Stat))
}

func (m *FileMock) Sync() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sync")
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *FileMockRecorder) Sync() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*FileMock)(nil).Sync))
}

func (m *FileMock) Truncate(size int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Truncate", size)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *FileMockRecorder) Truncate(size interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Truncate", reflect.TypeOf((*FileMock)(nil).Truncate), size)
}

func (m *FileMock) Write(p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *FileMockRecorder) Write(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*FileMock)(nil).Write), p)
}

func (m *FileMock) WriteAt(p []byte, off int64) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteAt", p, off)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *FileMockRecorder) WriteAt(p, off interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteAt", reflect.TypeOf((*FileMock)(nil).WriteAt), p, off)
}

func (m *FileMock) WriteString(s string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteString", s)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *FileMockRecorder) WriteString(s interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteString", reflect.TypeOf((*FileMock)(nil).WriteString), s)
}
