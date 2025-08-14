package tests

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
	config "github.com/happyhippyhippo/flam-config"
	mocks "github.com/happyhippyhippo/flam-config/tests/mocks"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
	time "github.com/happyhippyhippo/flam-time"
)

func Test_fileSource(t *testing.T) {
	t.Run("should ignore config without path field", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverFile,
				"disk":     "my_disk",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
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
				"driver":   config.SourceDriverFile,
				"disk":     "my_disk",
				"path":     "./testdata/invalid",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
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
				"driver":   config.SourceDriverFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
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
				"driver":   config.SourceDriverFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		expectedErr := errors.New("file error")
		disk := mocks.NewDisk(ctrl)
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
				"driver":   config.SourceDriverFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		expectedErr := errors.New("file error")
		file := mocks.NewFile(ctrl)
		file.EXPECT().Read(gomock.Any()).Return(0, expectedErr).Times(1)
		file.EXPECT().Close().Return(nil).Times(1)

		disk := mocks.NewDisk(ctrl)
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
				"driver":   config.SourceDriverFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		data := "{"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		file := mocks.NewFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)
		file.EXPECT().Close().Return(nil).Times(1)

		disk := mocks.NewDisk(ctrl)
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

	t.Run("should correctly load file source", func(t *testing.T) {
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
				"driver":   config.SourceDriverFile,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata/config",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		data := "field: value"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		file := mocks.NewFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)
		file.EXPECT().Close().Return(nil).Times(1)

		disk := mocks.NewDisk(ctrl)
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
			assert.NotNil(t, got)
			assert.NoError(t, e)

			assert.Equal(t, "value", got.Get("field"))
		}))
	})
}
