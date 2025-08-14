package mocks

import (
	os "os"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// File is a mock of FileSystemFile interface.
type File struct {
	ctrl     *gomock.Controller
	recorder *FileRecorder
}

// FileRecorder is the mock recorder for File.
type FileRecorder struct {
	mock *File
}

// NewFile creates a new mock instance.
func NewFile(ctrl *gomock.Controller) *File {
	mock := &File{ctrl: ctrl}
	mock.recorder = &FileRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *File) EXPECT() *FileRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *File) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *FileRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*File)(nil).Close))
}

// Name mocks base method.
func (m *File) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *FileRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*File)(nil).Name))
}

// Read mocks base method.
func (m *File) Read(p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *FileRecorder) Read(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*File)(nil).Read), p)
}

// ReadAt mocks base method.
func (m *File) ReadAt(p []byte, off int64) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadAt", p, off)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadAt indicates an expected call of ReadAt.
func (mr *FileRecorder) ReadAt(p, off interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAt", reflect.TypeOf((*File)(nil).ReadAt), p, off)
}

// Readdir mocks base method.
func (m *File) Readdir(count int) ([]os.FileInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Readdir", count)
	ret0, _ := ret[0].([]os.FileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Readdir indicates an expected call of Readdir.
func (mr *FileRecorder) Readdir(count interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Readdir", reflect.TypeOf((*File)(nil).Readdir), count)
}

// Readdirnames mocks base method.
func (m *File) Readdirnames(n int) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Readdirnames", n)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Readdirnames indicates an expected call of Readdirnames.
func (mr *FileRecorder) Readdirnames(n interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Readdirnames", reflect.TypeOf((*File)(nil).Readdirnames), n)
}

// Seek mocks base method.
func (m *File) Seek(offset int64, whence int) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Seek", offset, whence)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Seek indicates an expected call of Seek.
func (mr *FileRecorder) Seek(offset, whence interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Seek", reflect.TypeOf((*File)(nil).Seek), offset, whence)
}

// Stat mocks base method.
func (m *File) Stat() (os.FileInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stat")
	ret0, _ := ret[0].(os.FileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stat indicates an expected call of Stat.
func (mr *FileRecorder) Stat() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stat", reflect.TypeOf((*File)(nil).Stat))
}

// Sync mocks base method.
func (m *File) Sync() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sync")
	ret0, _ := ret[0].(error)
	return ret0
}

// Sync indicates an expected call of Sync.
func (mr *FileRecorder) Sync() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*File)(nil).Sync))
}

// Truncate mocks base method.
func (m *File) Truncate(size int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Truncate", size)
	ret0, _ := ret[0].(error)
	return ret0
}

// Truncate indicates an expected call of Truncate.
func (mr *FileRecorder) Truncate(size interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Truncate", reflect.TypeOf((*File)(nil).Truncate), size)
}

// Write mocks base method.
func (m *File) Write(p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *FileRecorder) Write(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*File)(nil).Write), p)
}

// WriteAt mocks base method.
func (m *File) WriteAt(p []byte, off int64) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteAt", p, off)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WriteAt indicates an expected call of WriteAt.
func (mr *FileRecorder) WriteAt(p, off interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteAt", reflect.TypeOf((*File)(nil).WriteAt), p, off)
}

// WriteString mocks base method.
func (m *File) WriteString(s string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteString", s)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WriteString indicates an expected call of WriteString.
func (mr *FileRecorder) WriteString(s interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteString", reflect.TypeOf((*File)(nil).WriteString), s)
}
