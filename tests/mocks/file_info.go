package mocks

import (
	fs "io/fs"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// FileInfo is a mock of FileInfo interface.
type FileInfo struct {
	ctrl     *gomock.Controller
	recorder *FileInfoRecorder
}

// FileInfoRecorder is the mock recorder for FileInfo.
type FileInfoRecorder struct {
	mock *FileInfo
}

// NewFileInfo creates a new mock instance.
func NewFileInfo(ctrl *gomock.Controller) *FileInfo {
	mock := &FileInfo{ctrl: ctrl}
	mock.recorder = &FileInfoRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *FileInfo) EXPECT() *FileInfoRecorder {
	return m.recorder
}

// IsDir mocks base method.
func (m *FileInfo) IsDir() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsDir")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsDir indicates an expected call of IsDir.
func (mr *FileInfoRecorder) IsDir() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsDir", reflect.TypeOf((*FileInfo)(nil).IsDir))
}

// ModTime mocks base method.
func (m *FileInfo) ModTime() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModTime")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// ModTime indicates an expected call of ModTime.
func (mr *FileInfoRecorder) ModTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModTime", reflect.TypeOf((*FileInfo)(nil).ModTime))
}

// Mode mocks base method.
func (m *FileInfo) Mode() fs.FileMode {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Mode")
	ret0, _ := ret[0].(fs.FileMode)
	return ret0
}

// Mode indicates an expected call of Mode.
func (mr *FileInfoRecorder) Mode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Mode", reflect.TypeOf((*FileInfo)(nil).Mode))
}

// Name mocks base method.
func (m *FileInfo) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *FileInfoRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*FileInfo)(nil).Name))
}

// Size mocks base method.
func (m *FileInfo) Size() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Size")
	ret0, _ := ret[0].(int64)
	return ret0
}

// Size indicates an expected call of Size.
func (mr *FileInfoRecorder) Size() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Size", reflect.TypeOf((*FileInfo)(nil).Size))
}

// Sys mocks base method.
func (m *FileInfo) Sys() interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sys")
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Sys indicates an expected call of Sys.
func (mr *FileInfoRecorder) Sys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sys", reflect.TypeOf((*FileInfo)(nil).Sys))
}
