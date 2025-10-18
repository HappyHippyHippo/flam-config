package config

import (
	"sync"
	"time"

	"go.uber.org/dig"

	"github.com/happyhippyhippo/flam"
	flamTime "github.com/happyhippyhippo/flam-time"
)

type provider struct {
	observer flamTime.Trigger
}

func NewProvider() flam.Provider {
	return &provider{}
}

func (*provider) Id() string {
	return providerId
}

func (*provider) Register(
	container *dig.Container,
) error {
	if container == nil {
		return newErrNilReference("container")
	}

	registerer := flam.NewRegisterer()
	registerer.Queue(newRestRequesterGenerator)
	registerer.Queue(newJsonParserCreator, dig.Group(ParserCreatorGroup))
	registerer.Queue(newYamlParserCreator, dig.Group(ParserCreatorGroup))
	registerer.Queue(newParserFactory)
	registerer.Queue(newEnvSourceCreator, dig.Group(SourceCreatorGroup))
	registerer.Queue(newFileSourceCreator, dig.Group(SourceCreatorGroup))
	registerer.Queue(newObservableFileSourceCreator, dig.Group(SourceCreatorGroup))
	registerer.Queue(newDirSourceCreator, dig.Group(SourceCreatorGroup))
	registerer.Queue(newRestSourceCreator, dig.Group(SourceCreatorGroup))
	registerer.Queue(newObservableRestSourceCreator, dig.Group(SourceCreatorGroup))
	registerer.Queue(newSourceFactory)
	registerer.Queue(newManager)
	registerer.Queue(newFactoryConfig)
	registerer.Queue(newFacade)

	return registerer.Run(container)
}

func (provider *provider) Boot(
	container *dig.Container,
) error {
	if container == nil {
		return newErrNilReference("container")
	}

	executor := flam.NewExecutor()
	executor.Queue(provider.bootDefaults)
	executor.Queue(provider.bootSources)
	executor.Queue(provider.bootObserver)

	return executor.Run(container)
}

func (provider *provider) Close(
	container *dig.Container,
) error {
	if container == nil {
		return newErrNilReference("container")
	}

	executor := flam.NewExecutor()
	executor.Queue(provider.closeObserver)
	executor.Queue(provider.closeSourceFactory)
	executor.Queue(provider.closeParserFactory)

	return executor.Run(container)
}

func (*provider) bootDefaults(
	manager *manager,
) error {
	defaultsSource := &source{mutex: &sync.Mutex{}, bag: Defaults, priority: -1}
	if e := manager.AddSource("defaults", defaultsSource); e != nil {
		return e
	}

	DefaultFileParser = manager.aggregate.String(PathDefaultFileParser, DefaultFileParser)
	DefaultFileDisk = manager.aggregate.String(PathDefaultFileDisk, DefaultFileDisk)
	DefaultRestParser = manager.aggregate.String(PathDefaultRestParser, DefaultRestParser)

	return nil
}

func (*provider) bootSources(
	manager *manager,
	sourceFactory sourceFactory,
) error {
	if manager.aggregate.Bool(PathBoot) {
		for id := range manager.aggregate.Bag(PathSources) {
			src, e := sourceFactory.Get(id)
			if e != nil {
				return e
			}

			if e = manager.AddSource(id, src); e != nil {
				return e
			}
		}
	}

	return nil
}

func (provider *provider) bootObserver(
	manager *manager,
	timeFacade flamTime.Facade,
) error {
	frequency := manager.aggregate.Duration(PathObserverFrequency)
	if frequency != time.Duration(0) {
		provider.observer, _ = timeFacade.NewRecurringTrigger(
			frequency,
			func() error {
				return manager.ReloadSources()
			})
	}

	return manager.AddObserver(
		"flam.config",
		PathObserverFrequency,
		func(old, new any) {
			newFrequency, ok := new.(time.Duration)
			if !ok {
				return
			}

			_ = provider.observer.Close()

			provider.observer, _ = timeFacade.NewRecurringTrigger(
				newFrequency,
				func() error {
					return manager.ReloadSources()
				})
		},
	)
}

func (provider *provider) closeObserver() error {
	if provider.observer == nil {
		return nil
	}

	return provider.observer.Close()
}

func (provider *provider) closeSourceFactory(
	sourceFactory sourceFactory,
) error {
	return sourceFactory.Close()
}

func (provider *provider) closeParserFactory(
	parserFactory parserFactory,
) error {
	return parserFactory.Close()
}
