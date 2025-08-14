package tests

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

	flam "github.com/happyhippyhippo/flam"
	config "github.com/happyhippyhippo/flam-config"
	mocks "github.com/happyhippyhippo/flam-config/tests/mocks"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
	flamTime "github.com/happyhippyhippo/flam-time"
)

func Test_NewProvider(t *testing.T) {
	assert.NotNil(t, config.NewProvider())
}

func Test_Provider_Id(t *testing.T) {
	assert.Equal(t, "flam.config.provider", config.NewProvider().Id())
}

func Test_Provider_Register(t *testing.T) {
	t.Run("should return error on nil container", func(t *testing.T) {
		assert.ErrorIs(
			t,
			config.NewProvider().Register(nil),
			flam.ErrNilReference)
	})

	t.Run("should successfully provide Facade", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.NotNil(t, facade)
		}))
	})

	t.Run("should successfully provide FactoryConfig", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(config flam.FactoryConfig) {
			assert.NotNil(t, config)
		}))
	})
}

func Test_Provider_Boot(t *testing.T) {
	t.Run("should return error on nil container", func(t *testing.T) {
		assert.ErrorIs(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(nil),
			flam.ErrNilReference)
	})

	t.Run("should return error on when trying to add the default source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		source := mocks.NewSource(ctrl)
		source.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{}).Times(1)
		require.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.NoError(t, facade.AddSource("defaults", source))
		}))

		assert.ErrorIs(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			config.ErrDuplicateSource)
	})

	t.Run("should correctly load the config defaults", func(t *testing.T) {
		config.Defaults = flam.Bag{"defaults": flam.Bag{"field": "value"}}
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.Equal(t, "value", facade.Get("defaults.field"))
		}))
	})

	t.Run("should use default boot values when not provided", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.Equal(t, "", config.DefaultFileDisk)
		assert.Equal(t, "", config.DefaultFileParser)
		assert.Equal(t, "", config.DefaultRestParser)
	})

	t.Run("should use provided default boot values when provided", func(t *testing.T) {
		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathDefaultFileDisk, "my_disk")
		_ = config.Defaults.Set(config.PathDefaultFileParser, "my_parser")
		_ = config.Defaults.Set(config.PathDefaultRestParser, "my_rest_parser")
		defer func() {
			config.DefaultFileDisk = ""
			config.DefaultFileParser = ""
			config.DefaultRestParser = ""
			config.Defaults = flam.Bag{}
		}()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.Equal(t, "my_disk", config.DefaultFileDisk)
		assert.Equal(t, "my_parser", config.DefaultFileParser)
		assert.Equal(t, "my_rest_parser", config.DefaultRestParser)
	})

	t.Run("should return source instantiation error", func(t *testing.T) {
		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   "invalid",
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

	t.Run("should return source storing error", func(t *testing.T) {
		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"defaults": flam.Bag{
				"driver":   config.SourceDriverEnv,
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
			config.ErrDuplicateSource)
	})

	t.Run("should correctly load the sources", func(t *testing.T) {
		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverEnv,
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.True(t, facade.HasSource("my_source"))
		}))
	})
}

func Test_Provider_Run(t *testing.T) {
	t.Run("should return error on nil container", func(t *testing.T) {
		assert.ErrorIs(t, config.NewProvider().(flam.RunnableProvider).Run(nil), flam.ErrNilReference)
	})

	t.Run("should register a config check frequency config observer", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		provider := config.NewProvider()
		require.NoError(t, provider.(flam.RunnableProvider).Run(container))
		defer func() { _ = provider.(flam.ClosableProvider).Close(container) }()

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.True(t, facade.HasObserver("flam.config", config.PathObserverFrequency))
		}))
	})

	t.Run("should update the config check frequency observer trigger when the config check frequency changes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_ = config.Defaults.Set(config.PathObserverFrequency, 10*time.Millisecond)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		timeFacade := mocks.NewTimeFacade(ctrl)
		trigger := mocks.NewTrigger(ctrl)
		gomock.InOrder(
			timeFacade.EXPECT().
				NewRecurringTrigger(10*time.Millisecond, gomock.Any()).
				Return(trigger, nil),
			trigger.EXPECT().Close().Return(nil),
			timeFacade.EXPECT().
				NewRecurringTrigger(20*time.Millisecond, gomock.Any()).
				Return(trigger, nil),
			trigger.EXPECT().Close().Return(nil),
		)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		provider := config.NewProvider()
		require.NoError(t, provider.(flam.RunnableProvider).Run(container))
		defer func() { _ = provider.(flam.ClosableProvider).Close(container) }()

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.NoError(t, facade.Set(config.PathObserverFrequency, 20*time.Millisecond))
		}))
	})

	t.Run("should not update the config check frequency observer trigger when the config check frequency is not a duration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_ = config.Defaults.Set(config.PathObserverFrequency, time.Second)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		timeFacade := mocks.NewTimeFacade(ctrl)
		trigger := mocks.NewTrigger(ctrl)
		trigger.EXPECT().Close().Return(nil)
		timeFacade.EXPECT().
			NewRecurringTrigger(time.Second, gomock.Any()).
			Return(trigger, nil).
			Times(1)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		provider := config.NewProvider()
		require.NoError(t, provider.(flam.RunnableProvider).Run(container))
		defer func() { _ = provider.(flam.ClosableProvider).Close(container) }()

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.NoError(t, facade.Set(config.PathObserverFrequency, "string"))
		}))
	})

	t.Run("should reload sources", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_ = config.Defaults.Set(config.PathObserverFrequency, 10*time.Millisecond)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		wg := &sync.WaitGroup{}
		source := mocks.NewObservableSource(ctrl)
		source.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{}).Times(3)
		source.EXPECT().GetPriority().Return(1).Times(1)
		source.EXPECT().Reload().DoAndReturn(func() (bool, error) {
			wg.Done()
			return false, nil
		}).Times(3)
		require.NoError(t, container.Invoke(func(facade config.Facade) error {
			return facade.AddSource("my_source", source)
		}))

		provider := config.NewProvider()
		require.NoError(t, provider.(flam.BootableProvider).Boot(container))

		wg.Add(1)
		require.NoError(t, provider.(flam.RunnableProvider).Run(container))
		defer func() { _ = provider.(flam.ClosableProvider).Close(container) }()
		wg.Wait() // load

		wg.Add(1)
		wg.Wait() // reload on trigger

		wg.Add(1)
		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.NoError(t, facade.Set(config.PathObserverFrequency, 20*time.Millisecond))
		}))
		wg.Wait() // reload on frequency change
	})
}

func Test_Provider_Close(t *testing.T) {
	t.Run("should return error on nil container", func(t *testing.T) {
		assert.ErrorIs(
			t,
			config.NewProvider().(flam.ClosableProvider).Close(nil),
			flam.ErrNilReference)
	})

	t.Run("should return parser closing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		provider := config.NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, provider.Register(container))

		expectedError := errors.New("close error")
		parser := mocks.NewParser(ctrl)
		parser.EXPECT().Close().Return(expectedError).Times(1)

		require.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.NoError(t, facade.AddParser("parser", parser))
		}))

		require.NoError(t, provider.(flam.BootableProvider).Boot(container))

		assert.ErrorIs(t, provider.(flam.ClosableProvider).Close(container), expectedError)
	})

	t.Run("should close a loaded json parser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(filesystem.PathDisks, flam.Bag{
			"my_disk": flam.Bag{
				"driver": "mock",
			}})
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"json": flam.Bag{
				"driver": config.ParserDriverJson,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"json": flam.Bag{
				"driver":   config.SourceDriverFile,
				"disk":     "my_disk",
				"path":     "config.json",
				"parser":   "json",
				"priority": 123,
			}})
		_ = config.Defaults.Set(config.PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		provider := config.NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, provider.Register(container))

		disk := afero.NewMemMapFs()
		json, _ := disk.Create("config.json")
		_, _ = json.Write([]byte(`{"key2": "value2"}`))

		diskCreatorConfig := flam.Bag{"id": "my_disk", "driver": "mock"}
		diskCreator := mocks.NewDiskCreator(ctrl)
		diskCreator.EXPECT().Accept(diskCreatorConfig).Return(true).Times(1)
		diskCreator.EXPECT().Create(diskCreatorConfig).Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.DiskCreator {
			return diskCreator
		}, dig.Group(filesystem.DiskCreatorGroup)))

		require.NoError(t, provider.(flam.BootableProvider).Boot(container))
		require.NoError(t, provider.(flam.RunnableProvider).Run(container))

		assert.NoError(t, provider.(flam.ClosableProvider).Close(container))
	})

	t.Run("should close a loaded yaml parser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(filesystem.PathDisks, flam.Bag{
			"my_disk": flam.Bag{
				"driver": "mock",
			}})
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"yaml": flam.Bag{
				"driver": config.ParserDriverYaml,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"yaml": flam.Bag{
				"driver":   config.SourceDriverFile,
				"disk":     "my_disk",
				"path":     "config.yaml",
				"parser":   "yaml",
				"priority": 123,
			}})
		_ = config.Defaults.Set(config.PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		provider := config.NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, provider.Register(container))

		disk := afero.NewMemMapFs()
		yaml, _ := disk.Create("config.yaml")
		_, _ = yaml.Write([]byte(`key1: value1`))

		diskCreatorConfig := flam.Bag{"id": "my_disk", "driver": "mock"}
		diskCreator := mocks.NewDiskCreator(ctrl)
		diskCreator.EXPECT().Accept(diskCreatorConfig).Return(true).Times(1)
		diskCreator.EXPECT().Create(diskCreatorConfig).Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.DiskCreator {
			return diskCreator
		}, dig.Group(filesystem.DiskCreatorGroup)))

		require.NoError(t, provider.(flam.BootableProvider).Boot(container))
		require.NoError(t, provider.(flam.RunnableProvider).Run(container))

		assert.NoError(t, provider.(flam.ClosableProvider).Close(container))
	})

	t.Run("should return source closing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver": "mock",
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		provider := config.NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, provider.Register(container))

		expectedError := errors.New("close error")
		source := mocks.NewSource(ctrl)
		source.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{}).Times(1)
		source.EXPECT().GetPriority().Return(1).Times(1)
		source.EXPECT().Close().Return(expectedError).Times(1)

		sourceCreatorConfig := flam.Bag{"id": "my_source", "driver": "mock"}
		sourceCreator := mocks.NewSourceCreator(ctrl)
		sourceCreator.EXPECT().Accept(sourceCreatorConfig).Return(true).Times(1)
		sourceCreator.EXPECT().Create(sourceCreatorConfig).Return(source, nil).Times(1)
		require.NoError(t, container.Provide(func() config.SourceCreator {
			return sourceCreator
		}, dig.Group(config.SourceCreatorGroup)))

		require.NoError(t, provider.(flam.BootableProvider).Boot(container))
		require.NoError(t, provider.(flam.RunnableProvider).Run(container))

		assert.ErrorIs(t, provider.(flam.ClosableProvider).Close(container), expectedError)
	})

	t.Run("should close the config observer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathObserverFrequency, 1000)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		provider := config.NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, provider.Register(container))

		require.NoError(t, provider.(flam.BootableProvider).Boot(container))
		require.NoError(t, provider.(flam.RunnableProvider).Run(container))

		assert.NoError(t, provider.(flam.ClosableProvider).Close(container))
	})
}
