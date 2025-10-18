package config

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"

	"github.com/happyhippyhippo/flam"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
	flamTime "github.com/happyhippyhippo/flam-time"
)

func Test_NewProvider(t *testing.T) {
	assert.NotNil(t, NewProvider())
}

func Test_Provider_Id(t *testing.T) {
	assert.Equal(t, "flam.config.provider", NewProvider().Id())
}

func Test_Provider_Register(t *testing.T) {
	t.Run("should return error on nil container", func(t *testing.T) {
		assert.ErrorIs(
			t,
			NewProvider().Register(nil),
			flam.ErrNilReference)
	})

	t.Run("should successfully provide Facade", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.NotNil(t, facade)
		}))
	})

	t.Run("should successfully provide FactoryConfig", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(config flam.FactoryConfig) {
			assert.NotNil(t, config)
		}))
	})
}

func Test_Provider_Boot(t *testing.T) {
	t.Run("should return error on nil container", func(t *testing.T) {
		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(nil),
			flam.ErrNilReference)
	})

	t.Run("should return error on when trying to add the default source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		src := NewSourceMock(ctrl)
		src.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{}).Times(1)
		require.NoError(t, container.Invoke(func(facade Facade) {
			assert.NoError(t, facade.AddSource("defaults", src))
		}))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			ErrDuplicateSource)
	})

	t.Run("should correctly load the config defaults", func(t *testing.T) {
		Defaults = flam.Bag{"defaults": flam.Bag{"field": "value"}}
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.Equal(t, "value", facade.Get("defaults.field"))
		}))
	})

	t.Run("should use default boot values when not provided", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.Equal(t, "", DefaultFileDisk)
		assert.Equal(t, "", DefaultFileParser)
		assert.Equal(t, "", DefaultRestParser)
	})

	t.Run("should use provided default boot values when provided", func(t *testing.T) {
		Defaults = flam.Bag{}
		_ = Defaults.Set(PathDefaultFileDisk, "my_disk")
		_ = Defaults.Set(PathDefaultFileParser, "my_parser")
		_ = Defaults.Set(PathDefaultRestParser, "my_rest_parser")
		defer func() {
			DefaultFileDisk = ""
			DefaultFileParser = ""
			DefaultRestParser = ""
			Defaults = flam.Bag{}
		}()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.Equal(t, "my_disk", DefaultFileDisk)
		assert.Equal(t, "my_parser", DefaultFileParser)
		assert.Equal(t, "my_rest_parser", DefaultRestParser)
	})

	t.Run("should return source instantiation error", func(t *testing.T) {
		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   "invalid",
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			flam.ErrInvalidResourceConfig)
	})

	t.Run("should return source storing error", func(t *testing.T) {
		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathSources, flam.Bag{
			"defaults": flam.Bag{
				"driver":   SourceDriverEnv,
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			ErrDuplicateSource)
	})

	t.Run("should correctly load the sources", func(t *testing.T) {
		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverEnv,
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.True(t, facade.HasSource("my_source"))
		}))
	})

	t.Run("should register a config check frequency config observer", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		p := NewProvider()
		_ = p.(flam.BootableProvider).Boot(container)
		defer func() { _ = p.(flam.ClosableProvider).Close(container) }()

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.True(t, facade.HasObserver("flam.config", PathObserverFrequency))
		}))
	})

	t.Run("should update the config check frequency observer trigger when the config check frequency changes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_ = Defaults.Set(PathObserverFrequency, 10*time.Millisecond)
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		timeFacade := NewTimeFacadeMock(ctrl)
		trigger := NewTriggerMock(ctrl)
		gomock.InOrder(
			timeFacade.EXPECT().
				NewRecurringTrigger(10*time.Millisecond, gomock.Any()).
				Return(trigger, nil),
			trigger.EXPECT().Close().Return(nil),
			timeFacade.EXPECT().
				NewRecurringTrigger(20*time.Millisecond, gomock.Any()).
				Return(trigger, nil),
		)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))

		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		p := NewProvider()
		defer func() { _ = p.(flam.ClosableProvider).Close(container) }()

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.NoError(t, facade.Set(PathObserverFrequency, 20*time.Millisecond))
		}))
	})

	t.Run("should not update the config check frequency observer trigger when the config check frequency is not a duration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_ = Defaults.Set(PathObserverFrequency, time.Second)
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		timeFacade := NewTimeFacadeMock(ctrl)
		trigger := NewTriggerMock(ctrl)
		timeFacade.EXPECT().
			NewRecurringTrigger(time.Second, gomock.Any()).
			Return(trigger, nil).
			Times(1)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))

		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		p := NewProvider()
		defer func() { _ = p.(flam.ClosableProvider).Close(container) }()

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.NoError(t, facade.Set(PathObserverFrequency, "string"))
		}))
	})

	t.Run("should reload sources", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_ = Defaults.Set(PathObserverFrequency, 10*time.Millisecond)
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		wg := &sync.WaitGroup{}
		src := NewObservableSourceMock(ctrl)
		src.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{}).Times(3)
		src.EXPECT().GetPriority().Return(1).Times(1)
		src.EXPECT().Reload().DoAndReturn(func() (bool, error) {
			wg.Done()
			return false, nil
		}).Times(3)
		require.NoError(t, container.Invoke(func(facade Facade) error {
			return facade.AddSource("my_source", src)
		}))

		p := NewProvider()
		require.NoError(t, p.(flam.BootableProvider).Boot(container))

		wg.Add(1)
		defer func() { _ = p.(flam.ClosableProvider).Close(container) }()
		wg.Wait() // load

		wg.Add(1)
		wg.Wait() // reload on trigger

		wg.Add(1)
		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.NoError(t, facade.Set(PathObserverFrequency, 20*time.Millisecond))
		}))
		wg.Wait() // reload on frequency change
	})
}

func Test_Provider_Close(t *testing.T) {
	t.Run("should return error on nil container", func(t *testing.T) {
		assert.ErrorIs(
			t,
			NewProvider().(flam.ClosableProvider).Close(nil),
			flam.ErrNilReference)
	})

	t.Run("should return parser closing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		p := NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, p.Register(container))

		expectedError := errors.New("close error")
		parser := NewParserMock(ctrl)
		parser.EXPECT().Close().Return(expectedError).Times(1)

		require.NoError(t, container.Invoke(func(facade Facade) {
			assert.NoError(t, facade.AddParser("parser", parser))
		}))

		require.NoError(t, p.(flam.BootableProvider).Boot(container))

		assert.ErrorIs(t, p.(flam.ClosableProvider).Close(container), expectedError)
	})

	t.Run("should close a loaded json parser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(filesystem.PathDisks, flam.Bag{
			"my_disk": flam.Bag{
				"driver": "mock",
			}})
		_ = Defaults.Set(PathParsers, flam.Bag{
			"json": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"json": flam.Bag{
				"driver":   SourceDriverFile,
				"disk":     "my_disk",
				"path":     "config.json",
				"parser":   "json",
				"priority": 123,
			}})
		_ = Defaults.Set(PathBoot, true)
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		p := NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, p.Register(container))

		disk := afero.NewMemMapFs()
		json, _ := disk.Create("config.json")
		_, _ = json.Write([]byte(`{"key2": "value2"}`))

		diskCreatorConfig := flam.Bag{"id": "my_disk", "driver": "mock"}
		diskCreator := NewDiskCreatorMock(ctrl)
		diskCreator.EXPECT().Accept(diskCreatorConfig).Return(true).Times(1)
		diskCreator.EXPECT().Create(diskCreatorConfig).Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.DiskCreator {
			return diskCreator
		}, dig.Group(filesystem.DiskCreatorGroup)))

		require.NoError(t, p.(flam.BootableProvider).Boot(container))

		assert.NoError(t, p.(flam.ClosableProvider).Close(container))
	})

	t.Run("should close a loaded yaml parser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(filesystem.PathDisks, flam.Bag{
			"my_disk": flam.Bag{
				"driver": "mock",
			}})
		_ = Defaults.Set(PathParsers, flam.Bag{
			"yaml": flam.Bag{
				"driver": ParserDriverYaml,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"yaml": flam.Bag{
				"driver":   SourceDriverFile,
				"disk":     "my_disk",
				"path":     "config.yaml",
				"parser":   "yaml",
				"priority": 123,
			}})
		_ = Defaults.Set(PathBoot, true)
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		p := NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, p.Register(container))

		disk := afero.NewMemMapFs()
		yaml, _ := disk.Create("config.yaml")
		_, _ = yaml.Write([]byte(`key1: value1`))

		diskCreatorConfig := flam.Bag{"id": "my_disk", "driver": "mock"}
		diskCreator := NewDiskCreatorMock(ctrl)
		diskCreator.EXPECT().Accept(diskCreatorConfig).Return(true).Times(1)
		diskCreator.EXPECT().Create(diskCreatorConfig).Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.DiskCreator {
			return diskCreator
		}, dig.Group(filesystem.DiskCreatorGroup)))

		require.NoError(t, p.(flam.BootableProvider).Boot(container))

		assert.NoError(t, p.(flam.ClosableProvider).Close(container))
	})

	t.Run("should return source closing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver": "mock",
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		p := NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, p.Register(container))

		expectedError := errors.New("close error")
		src := NewSourceMock(ctrl)
		src.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{}).Times(1)
		src.EXPECT().GetPriority().Return(1).Times(1)
		src.EXPECT().Close().Return(expectedError).Times(1)

		sourceCreatorConfig := flam.Bag{"id": "my_source", "driver": "mock"}
		sourceCreator := NewSourceCreatorMock(ctrl)
		sourceCreator.EXPECT().Accept(sourceCreatorConfig).Return(true).Times(1)
		sourceCreator.EXPECT().Create(sourceCreatorConfig).Return(src, nil).Times(1)
		require.NoError(t, container.Provide(func() SourceCreator {
			return sourceCreator
		}, dig.Group(SourceCreatorGroup)))

		require.NoError(t, p.(flam.BootableProvider).Boot(container))

		assert.ErrorIs(t, p.(flam.ClosableProvider).Close(container), expectedError)
	})

	t.Run("should close the config observer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathObserverFrequency, 1000)
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		p := NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, p.Register(container))

		require.NoError(t, p.(flam.BootableProvider).Boot(container))

		assert.NoError(t, p.(flam.ClosableProvider).Close(container))
	})
}
