package config

import (
	fs "io/fs"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

type FileInfoMock struct {
	ctrl     *gomock.Controller
	recorder *FileInfoMockRecorder
}

type FileInfoMockRecorder struct {
	mock *FileInfoMock
}

func NewFileInfoMock(ctrl *gomock.Controller) *FileInfoMock {
	mock := &FileInfoMock{ctrl: ctrl}
	mock.recorder = &FileInfoMockRecorder{mock}
	return mock
}

func (m *FileInfoMock) EXPECT() *FileInfoMockRecorder {
	return m.recorder
}

func (m *FileInfoMock) IsDir() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsDir")
	ret0, _ := ret[0].(bool)
	return ret0
}

func (mr *FileInfoMockRecorder) IsDir() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsDir", reflect.TypeOf((*FileInfoMock)(nil).IsDir))
}

func (m *FileInfoMock) ModTime() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModTime")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

func (mr *FileInfoMockRecorder) ModTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModTime", reflect.TypeOf((*FileInfoMock)(nil).ModTime))
}

func (m *FileInfoMock) Mode() fs.FileMode {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Mode")
	ret0, _ := ret[0].(fs.FileMode)
	return ret0
}

func (mr *FileInfoMockRecorder) Mode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Mode", reflect.TypeOf((*FileInfoMock)(nil).Mode))
}

func (m *FileInfoMock) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

func (mr *FileInfoMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*FileInfoMock)(nil).Name))
}

func (m *FileInfoMock) Size() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Size")
	ret0, _ := ret[0].(int64)
	return ret0
}

func (mr *FileInfoMockRecorder) Size() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Size", reflect.TypeOf((*FileInfoMock)(nil).Size))
}

func (m *FileInfoMock) Sys() interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sys")
	ret0 := ret[0]
	return ret0
}

func (mr *FileInfoMockRecorder) Sys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sys", reflect.TypeOf((*FileInfoMock)(nil).Sys))
}
