package tests

import (
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
	config "github.com/happyhippyhippo/flam-config"
	mocks "github.com/happyhippyhippo/flam-config/tests/mocks"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
	flamTime "github.com/happyhippyhippo/flam-time"
)

func Test_observableFileSource(t *testing.T) {
	t.Run("should ignore config without path field", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverObservableFile,
				"disk":     "my_disk",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		assert.ErrorIs(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			flam.ErrInvalidResourceConfig)
	})

	t.Run("should return filesystem disk retrieval error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverObservableFile,
				"disk":     "my_disk",
				"path":     "./testdata/invalid",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		expectedErr := errors.New("filesystem error")
		fsFacade := mocks.NewFileSystemFacade(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(nil, expectedErr).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		assert.ErrorIs(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			expectedErr)
	})

	t.Run("should return parser retrieval error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverObservableFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		disk := mocks.NewDisk(ctrl)

		fsFacade := mocks.NewFileSystemFacade(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		assert.ErrorIs(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			flam.ErrUnknownResource)
	})

	t.Run("should return file stat error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": config.ParserDriverYaml,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverObservableFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		expectedErr := errors.New("file error")
		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Stat("/testdata/config").Return(nil, expectedErr).Times(1)

		fsFacade := mocks.NewFileSystemFacade(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		assert.ErrorIs(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			expectedErr)
	})

	t.Run("should return file opening error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": config.ParserDriverYaml,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverObservableFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		now := time.Now()
		fileInfo := mocks.NewFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(now).Times(1)

		expectedErr := errors.New("file error")
		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Stat("/testdata/config").Return(fileInfo, nil).Times(1)
		disk.EXPECT().
			OpenFile("/testdata/config", os.O_RDONLY, os.FileMode(0o644)).
			Return(nil, expectedErr).
			Times(1)

		fsFacade := mocks.NewFileSystemFacade(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		assert.ErrorIs(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			expectedErr)
	})

	t.Run("should return file reading error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": config.ParserDriverYaml,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverObservableFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		expectedErr := errors.New("file error")
		file := mocks.NewFile(ctrl)
		file.EXPECT().Read(gomock.Any()).Return(0, expectedErr).Times(1)
		file.EXPECT().Close().Return(nil).Times(1)

		now := time.Now()
		fileInfo := mocks.NewFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(now).Times(1)

		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Stat("/testdata/config").Return(fileInfo, nil).Times(1)
		disk.EXPECT().
			OpenFile("/testdata/config", os.O_RDONLY, os.FileMode(0o644)).
			Return(file, nil).
			Times(1)

		fsFacade := mocks.NewFileSystemFacade(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		assert.ErrorIs(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			expectedErr)
	})

	t.Run("should return file parsing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": config.ParserDriverYaml,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverObservableFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		data := "{"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		file := mocks.NewFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)
		file.EXPECT().Close().Return(nil).Times(1)

		now := time.Now()
		fileInfo := mocks.NewFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(now).Times(1)

		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Stat("/testdata/config").Return(fileInfo, nil).Times(1)
		disk.EXPECT().
			OpenFile("/testdata/config", os.O_RDONLY, os.FileMode(0o644)).
			Return(file, nil).
			Times(1)

		fsFacade := mocks.NewFileSystemFacade(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		assert.ErrorContains(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			"yaml: line 1: did not find expected node content")
	})

	t.Run("should correctly load observable file source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": config.ParserDriverYaml,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverObservableFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		data := "field: value"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		file := mocks.NewFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)
		file.EXPECT().Close().Return(nil).Times(1)

		now := time.Now()
		fileInfo := mocks.NewFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(now).Times(1)

		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Stat("/testdata/config").Return(fileInfo, nil).Times(1)
		disk.EXPECT().
			OpenFile("/testdata/config", os.O_RDONLY, os.FileMode(0o644)).
			Return(file, nil).
			Times(1)

		fsFacade := mocks.NewFileSystemFacade(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			assert.Equal(t, "value", got.Get("field"))
		}))
	})
}

func Test_observableFileSource_Reload(t *testing.T) {
	t.Run("should return file stat error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": config.ParserDriverYaml,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverObservableFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		data := "field: value"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		file := mocks.NewFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)
		file.EXPECT().Close().Return(nil).Times(1)

		now := time.Now()
		fileInfo := mocks.NewFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(now).Times(1)

		expectedErr := errors.New("filesystem error")
		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Stat("/testdata/config").Return(fileInfo, nil)
		disk.EXPECT().Stat("/testdata/config").Return(nil, expectedErr)
		disk.EXPECT().
			OpenFile("/testdata/config", os.O_RDONLY, os.FileMode(0o644)).
			Return(file, nil).
			Times(1)

		fsFacade := mocks.NewFileSystemFacade(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			reloaded, e := got.(config.ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorIs(t, e, expectedErr)
		}))
	})

	t.Run("should no-op if the file was not updated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": config.ParserDriverYaml,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverObservableFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		data := "field: value"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		file := mocks.NewFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)
		file.EXPECT().Close().Return(nil).Times(1)

		now := time.Now()
		fileInfo := mocks.NewFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(now).Times(2)

		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Stat("/testdata/config").Return(fileInfo, nil).Times(2)
		disk.EXPECT().
			OpenFile("/testdata/config", os.O_RDONLY, os.FileMode(0o644)).
			Return(file, nil).
			Times(1)

		fsFacade := mocks.NewFileSystemFacade(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			reloaded, e := got.(config.ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.NoError(t, e)
		}))
	})

	t.Run("should return file opening error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": config.ParserDriverYaml,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverObservableFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		data := "field: value"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		file := mocks.NewFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)
		file.EXPECT().Close().Return(nil).Times(1)

		now := time.Now()
		future := now.AddDate(1, 0, 0)
		fileInfo := mocks.NewFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(now)
		fileInfo.EXPECT().ModTime().Return(future)

		expectedErr := errors.New("filesystem error")
		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Stat("/testdata/config").Return(fileInfo, nil).Times(2)
		disk.EXPECT().
			OpenFile("/testdata/config", os.O_RDONLY, os.FileMode(0o644)).
			Return(file, nil).
			Times(1)
		disk.EXPECT().
			OpenFile("/testdata/config", os.O_RDONLY, os.FileMode(0o644)).
			Return(nil, expectedErr)

		fsFacade := mocks.NewFileSystemFacade(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			reloaded, e := got.(config.ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorIs(t, e, expectedErr)
		}))
	})

	t.Run("should return file reading error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": config.ParserDriverYaml,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverObservableFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		data := "field: value"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		expectedErr := errors.New("filesystem error")
		file := mocks.NewFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(reader)
		file.EXPECT().Read(gomock.Any()).Return(0, expectedErr)
		file.EXPECT().Close().Return(nil).Times(2)

		now := time.Now()
		future := now.AddDate(1, 0, 0)
		fileInfo := mocks.NewFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(now)
		fileInfo.EXPECT().ModTime().Return(future)

		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Stat("/testdata/config").Return(fileInfo, nil).Times(2)
		disk.EXPECT().
			OpenFile("/testdata/config", os.O_RDONLY, os.FileMode(0o644)).
			Return(file, nil).
			Times(2)

		fsFacade := mocks.NewFileSystemFacade(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			reloaded, e := got.(config.ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorIs(t, e, expectedErr)
		}))
	})

	t.Run("should return file parsing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": config.ParserDriverYaml,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverObservableFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		data1 := "field: value"
		reader1 := func(b []byte) (int, error) {
			copy(b, data1)
			return len(data1), io.EOF
		}

		data2 := "{"
		reader2 := func(b []byte) (int, error) {
			copy(b, data2)
			return len(data2), io.EOF
		}

		file := mocks.NewFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(reader1)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(reader2)
		file.EXPECT().Close().Return(nil).Times(2)

		now := time.Now()
		future := now.AddDate(1, 0, 0)
		fileInfo := mocks.NewFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(now)
		fileInfo.EXPECT().ModTime().Return(future)

		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Stat("/testdata/config").Return(fileInfo, nil).Times(2)
		disk.EXPECT().
			OpenFile("/testdata/config", os.O_RDONLY, os.FileMode(0o644)).
			Return(file, nil).
			Times(2)

		fsFacade := mocks.NewFileSystemFacade(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			reloaded, e := got.(config.ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorContains(t, e, "yaml: line 1: did not find expected node content")
		}))
	})

	t.Run("should correctly reload the source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": config.ParserDriverYaml,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverObservableFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		data1 := "field: value"
		reader1 := func(b []byte) (int, error) {
			copy(b, data1)
			return len(data1), io.EOF
		}

		data2 := "field2: value2"
		reader2 := func(b []byte) (int, error) {
			copy(b, data2)
			return len(data2), io.EOF
		}

		file := mocks.NewFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(reader1)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(reader2)
		file.EXPECT().Close().Return(nil).Times(2)

		now := time.Now()
		future := now.AddDate(1, 0, 0)
		fileInfo := mocks.NewFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(now)
		fileInfo.EXPECT().ModTime().Return(future)

		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Stat("/testdata/config").Return(fileInfo, nil).Times(2)
		disk.EXPECT().
			OpenFile("/testdata/config", os.O_RDONLY, os.FileMode(0o644)).
			Return(file, nil).
			Times(2)

		fsFacade := mocks.NewFileSystemFacade(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			reloaded, e := got.(config.ObservableSource).Reload()
			require.True(t, reloaded)
			require.NoError(t, e)

			assert.Nil(t, got.Get("field"))
			assert.Equal(t, "value2", got.Get("field2"))
		}))
	})
}
