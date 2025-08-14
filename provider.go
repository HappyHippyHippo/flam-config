package config

import (
	"sync"
	"time"

	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
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

	var e error
	provide := func(constructor any, opts ...dig.ProvideOption) bool {
		e = container.Provide(constructor, opts...)
		return e == nil
	}

	_ = provide(newRestRequesterGenerator) &&
		provide(newJsonParserCreator, dig.Group(ParserCreatorGroup)) &&
		provide(newYamlParserCreator, dig.Group(ParserCreatorGroup)) &&
		provide(newParserFactory) &&
		provide(newEnvSourceCreator, dig.Group(SourceCreatorGroup)) &&
		provide(newFileSourceCreator, dig.Group(SourceCreatorGroup)) &&
		provide(newObservableFileSourceCreator, dig.Group(SourceCreatorGroup)) &&
		provide(newDirSourceCreator, dig.Group(SourceCreatorGroup)) &&
		provide(newRestSourceCreator, dig.Group(SourceCreatorGroup)) &&
		provide(newObservableRestSourceCreator, dig.Group(SourceCreatorGroup)) &&
		provide(newSourceFactory) &&
		provide(newManager) &&
		provide(newFactoryConfig) &&
		provide(newFacade)

	return e
}

func (provider *provider) Boot(
	container *dig.Container,
) error {
	if container == nil {
		return newErrNilReference("container")
	}

	return container.Invoke(func(
		manager *manager,
		sourceFactory sourceFactory,
	) error {
		defaultsSource := &source{mutex: &sync.Mutex{}, bag: Defaults, priority: -1}
		if e := manager.AddSource("defaults", defaultsSource); e != nil {
			return e
		}

		DefaultFileParser = manager.aggregate.String(PathDefaultFileParser, DefaultFileParser)
		DefaultFileDisk = manager.aggregate.String(PathDefaultFileDisk, DefaultFileDisk)
		DefaultRestParser = manager.aggregate.String(PathDefaultRestParser, DefaultRestParser)

		if manager.aggregate.Bool(PathBoot) {
			for id := range manager.aggregate.Bag(PathSources) {
				source, e := sourceFactory.Get(id)
				if e != nil {
					return e
				}

				if e = manager.AddSource(id, source); e != nil {
					return e
				}
			}
		}

		return nil
	})
}

func (provider *provider) Run(
	container *dig.Container,
) error {
	if container == nil {
		return newErrNilReference("container")
	}

	return container.Invoke(func(
		manager *manager,
		timeFacade flamTime.Facade,
	) error {
		frequency := manager.aggregate.Duration(PathObserverFrequency)
		if frequency != time.Duration(0) {
			provider.observer, _ = timeFacade.NewRecurringTrigger(frequency, func() error {
				return manager.ReloadSources()
			})
		}

		return manager.AddObserver(
			"flam.config",
			PathObserverFrequency,
			func(old, new any) {
				frequency, ok := new.(time.Duration)
				if !ok {
					return
				}

				_ = provider.observer.Close()

				provider.observer, _ = timeFacade.NewRecurringTrigger(frequency, func() error {
					return manager.ReloadSources()
				})
			},
		)
	})
}

func (provider *provider) Close(
	container *dig.Container,
) error {
	if container == nil {
		return newErrNilReference("container")
	}

	return container.Invoke(func(
		sourceFactory sourceFactory,
		parserFactory parserFactory,
	) error {
		if provider.observer != nil {
			_ = provider.observer.Close()
		}

		if e := sourceFactory.Close(); e != nil {
			return e
		}

		if e := parserFactory.Close(); e != nil {
			return e
		}

		return nil
	})
}
