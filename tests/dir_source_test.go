package tests

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
	config "github.com/happyhippyhippo/flam-config"
	mocks "github.com/happyhippyhippo/flam-config/tests/mocks"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
	time "github.com/happyhippyhippo/flam-time"
)

func Test_dirSource(t *testing.T) {
	t.Run("should ignore config without path field", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverDir,
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
				"driver":   config.SourceDriverDir,
				"disk":     "my_disk",
				"path":     "./testdata",
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
				"driver":   config.SourceDriverDir,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		disk := afero.NewMemMapFs()

		fsFacade := mocks.NewFileSystemFacade(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		assert.ErrorIs(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			flam.ErrUnknownResource)
	})

	t.Run("should return dir opening error", func(t *testing.T) {
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
				"driver":   config.SourceDriverDir,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		disk := afero.NewMemMapFs()

		fsFacade := mocks.NewFileSystemFacade(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		assert.ErrorContains(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			"file does not exist")
	})

	t.Run("should return dir reading error", func(t *testing.T) {
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
				"driver":   config.SourceDriverDir,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		expectedErr := errors.New("dir error")
		dir := mocks.NewFile(ctrl)
		dir.EXPECT().Readdir(0).Return(nil, expectedErr).Times(1)
		dir.EXPECT().Close().Return(nil).Times(1)

		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Open("/testdata").Return(dir, nil).Times(1)

		fsFacade := mocks.NewFileSystemFacade(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		assert.ErrorIs(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			expectedErr)
	})

	t.Run("should correctly load an empty directory", func(t *testing.T) {
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
				"driver":   config.SourceDriverDir,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		dir := mocks.NewFile(ctrl)
		dir.EXPECT().Readdir(0).Return([]os.FileInfo{}, nil).Times(1)
		dir.EXPECT().Close().Return(nil).Times(1)

		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Open("/testdata").Return(dir, nil).Times(1)

		fsFacade := mocks.NewFileSystemFacade(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, container.Invoke(func(facade config.Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			bag := got.Get("", nil)
			assert.Empty(t, bag)
		}))
	})

	t.Run("should return the directory file opening error", func(t *testing.T) {
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
				"driver":   config.SourceDriverDir,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		fileInfo := mocks.NewFileInfo(ctrl)
		fileInfo.EXPECT().IsDir().Return(false).Times(1)
		fileInfo.EXPECT().Name().Return("file.yaml").Times(1)

		dir := mocks.NewFile(ctrl)
		dir.EXPECT().Readdir(0).Return([]os.FileInfo{fileInfo}, nil).Times(1)
		dir.EXPECT().Close().Return(nil).Times(1)

		expectedErr := errors.New("file error")
		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Open("/testdata").Return(dir, nil).Times(1)
		disk.EXPECT().
			OpenFile("/testdata/file.yaml", os.O_RDONLY, os.FileMode(0o644)).
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

	t.Run("should return the directory file parsing error", func(t *testing.T) {
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
				"driver":   config.SourceDriverDir,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata",
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

		fileInfo := mocks.NewFileInfo(ctrl)
		fileInfo.EXPECT().IsDir().Return(false).Times(1)
		fileInfo.EXPECT().Name().Return("file.yaml").Times(1)

		dir := mocks.NewFile(ctrl)
		dir.EXPECT().Readdir(0).Return([]os.FileInfo{fileInfo}, nil).Times(1)
		dir.EXPECT().Close().Return(nil).Times(1)

		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Open("/testdata").Return(dir, nil).Times(1)
		disk.EXPECT().
			OpenFile("/testdata/file.yaml", os.O_RDONLY, os.FileMode(0o644)).
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

	t.Run("should correctly load directory files", func(t *testing.T) {
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
				"driver":   config.SourceDriverDir,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata",
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

		fileInfo := mocks.NewFileInfo(ctrl)
		fileInfo.EXPECT().IsDir().Return(false).Times(1)
		fileInfo.EXPECT().Name().Return("file.yaml").Times(1)

		dir := mocks.NewFile(ctrl)
		dir.EXPECT().Readdir(0).Return([]os.FileInfo{fileInfo}, nil).Times(1)
		dir.EXPECT().Close().Return(nil).Times(1)

		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Open("/testdata").Return(dir, nil).Times(1)
		disk.EXPECT().
			OpenFile("/testdata/file.yaml", os.O_RDONLY, os.FileMode(0o644)).
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

	t.Run("should not load sub-directories if not flagged as recursive", func(t *testing.T) {
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
				"driver":   config.SourceDriverDir,
				"disk":     "my_disk",
				"parser":   "my_parser",
				"path":     "/testdata",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		subDirInfo := mocks.NewFileInfo(ctrl)
		subDirInfo.EXPECT().IsDir().Return(true).Times(1)

		dir := mocks.NewFile(ctrl)
		dir.EXPECT().Readdir(0).Return([]os.FileInfo{subDirInfo}, nil).Times(1)
		dir.EXPECT().Close().Return(nil).Times(1)

		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Open("/testdata").Return(dir, nil).Times(1)

		fsFacade := mocks.NewFileSystemFacade(ctrl)
		fsFacade.EXPECT().GetDisk("my_disk").Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.Facade { return fsFacade }))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			assert.Nil(t, got.Get("field"))
		}))
	})

	t.Run("should return the sub-directory files opening error", func(t *testing.T) {
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
				"driver":    config.SourceDriverDir,
				"disk":      "my_disk",
				"parser":    "my_parser",
				"path":      "/testdata",
				"recursive": true,
				"priority":  123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		fileInfo := mocks.NewFileInfo(ctrl)
		fileInfo.EXPECT().IsDir().Return(false).Times(1)
		fileInfo.EXPECT().Name().Return("file.yaml").Times(1)

		subDir := mocks.NewFile(ctrl)
		subDir.EXPECT().Readdir(0).Return([]os.FileInfo{fileInfo}, nil).Times(1)
		subDir.EXPECT().Close().Return(nil).Times(1)

		subDirInfo := mocks.NewFileInfo(ctrl)
		subDirInfo.EXPECT().IsDir().Return(true).Times(1)
		subDirInfo.EXPECT().Name().Return("subdir").Times(1)

		dir := mocks.NewFile(ctrl)
		dir.EXPECT().Readdir(0).Return([]os.FileInfo{subDirInfo}, nil).Times(1)
		dir.EXPECT().Close().Return(nil).Times(1)

		expectedErr := errors.New("file error")
		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Open("/testdata").Return(dir, nil).Times(1)
		disk.EXPECT().Open("/testdata/subdir").Return(subDir, nil).Times(1)
		disk.EXPECT().
			OpenFile("/testdata/subdir/file.yaml", os.O_RDONLY, os.FileMode(0o644)).
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

	t.Run("should return the sub-directory files parsing error", func(t *testing.T) {
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
				"driver":    config.SourceDriverDir,
				"disk":      "my_disk",
				"parser":    "my_parser",
				"path":      "/testdata",
				"recursive": true,
				"priority":  123,
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

		fileInfo := mocks.NewFileInfo(ctrl)
		fileInfo.EXPECT().IsDir().Return(false).Times(1)
		fileInfo.EXPECT().Name().Return("file.yaml").Times(1)

		subDir := mocks.NewFile(ctrl)
		subDir.EXPECT().Readdir(0).Return([]os.FileInfo{fileInfo}, nil).Times(1)
		subDir.EXPECT().Close().Return(nil).Times(1)

		subDirInfo := mocks.NewFileInfo(ctrl)
		subDirInfo.EXPECT().IsDir().Return(true).Times(1)
		subDirInfo.EXPECT().Name().Return("subdir").Times(1)

		dir := mocks.NewFile(ctrl)
		dir.EXPECT().Readdir(0).Return([]os.FileInfo{subDirInfo}, nil).Times(1)
		dir.EXPECT().Close().Return(nil).Times(1)

		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Open("/testdata").Return(dir, nil).Times(1)
		disk.EXPECT().Open("/testdata/subdir").Return(subDir, nil).Times(1)
		disk.EXPECT().
			OpenFile("/testdata/subdir/file.yaml", os.O_RDONLY, os.FileMode(0o644)).
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

	t.Run("should correctly load sub-directory files if flagged as recursive", func(t *testing.T) {
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
				"driver":    config.SourceDriverDir,
				"disk":      "my_disk",
				"parser":    "my_parser",
				"path":      "/testdata",
				"recursive": true,
				"priority":  123,
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

		fileInfo := mocks.NewFileInfo(ctrl)
		fileInfo.EXPECT().IsDir().Return(false).Times(1)
		fileInfo.EXPECT().Name().Return("file.yaml").Times(1)

		subDir := mocks.NewFile(ctrl)
		subDir.EXPECT().Readdir(0).Return([]os.FileInfo{fileInfo}, nil).Times(1)
		subDir.EXPECT().Close().Return(nil).Times(1)

		subDirInfo := mocks.NewFileInfo(ctrl)
		subDirInfo.EXPECT().IsDir().Return(true).Times(1)
		subDirInfo.EXPECT().Name().Return("subdir").Times(1)

		dir := mocks.NewFile(ctrl)
		dir.EXPECT().Readdir(0).Return([]os.FileInfo{subDirInfo}, nil).Times(1)
		dir.EXPECT().Close().Return(nil).Times(1)

		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Open("/testdata").Return(dir, nil).Times(1)
		disk.EXPECT().Open("/testdata/subdir").Return(subDir, nil).Times(1)
		disk.EXPECT().
			OpenFile("/testdata/subdir/file.yaml", os.O_RDONLY, os.FileMode(0o644)).
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
