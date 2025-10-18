package config

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
	time "github.com/happyhippyhippo/flam-time"
)

func Test_fileSource(t *testing.T) {
	t.Run("should ignore config without path field", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverFile,
				"disk":     "my_disk",
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			flam.ErrInvalidResourceConfig)
	})

	t.Run("should return filesystem disk retrieval error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverFile,
				"disk":     "my_disk",
				"path":     "./testdata/invalid",
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		expectedErr := errors.New("filesystem error")
		fsFacade := NewFileSystemFacadeMock(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(nil, expectedErr).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			expectedErr)
	})

	t.Run("should return parser retrieval error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		disk := NewDiskMock(ctrl)

		fsFacade := NewFileSystemFacadeMock(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			flam.ErrUnknownResource)
	})

	t.Run("should return file opening error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverYaml,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		expectedErr := errors.New("file error")
		disk := NewDiskMock(ctrl)
		disk.EXPECT().
			OpenFile("/testdata/config", os.O_RDONLY, os.FileMode(0o644)).
			Return(nil, expectedErr).
			Times(1)

		fsFacade := NewFileSystemFacadeMock(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			expectedErr)
	})

	t.Run("should return file reading error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverYaml,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		expectedErr := errors.New("file error")
		file := NewFileMock(ctrl)
		file.EXPECT().Read(gomock.Any()).Return(0, expectedErr).Times(1)
		file.EXPECT().Close().Return(nil).Times(1)

		disk := NewDiskMock(ctrl)
		disk.EXPECT().
			OpenFile("/testdata/config", os.O_RDONLY, os.FileMode(0o644)).
			Return(file, nil).
			Times(1)

		fsFacade := NewFileSystemFacadeMock(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			expectedErr)
	})

	t.Run("should return file parsing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverYaml,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		data := "{"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		file := NewFileMock(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)
		file.EXPECT().Close().Return(nil).Times(1)

		disk := NewDiskMock(ctrl)
		disk.EXPECT().
			OpenFile("/testdata/config", os.O_RDONLY, os.FileMode(0o644)).
			Return(file, nil).
			Times(1)

		fsFacade := NewFileSystemFacadeMock(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		assert.ErrorContains(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			"yaml: line 1: did not find expected node content")
	})

	t.Run("should correctly load file source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverYaml,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		data := "field: value"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		file := NewFileMock(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)
		file.EXPECT().Close().Return(nil).Times(1)

		disk := NewDiskMock(ctrl)
		disk.EXPECT().
			OpenFile("/testdata/config", os.O_RDONLY, os.FileMode(0o644)).
			Return(file, nil).
			Times(1)

		fsFacade := NewFileSystemFacadeMock(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.GetSource("my_source")
			assert.NotNil(t, got)
			assert.NoError(t, e)

			assert.Equal(t, "value", got.Get("field"))
		}))
	})
}
