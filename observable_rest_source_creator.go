package config

import (
	flam "github.com/happyhippyhippo/flam"
	time "github.com/happyhippyhippo/flam-time"
)

type observableRestSourceCreator struct {
	restRequesterGenerator RestRequesterGenerator
	parserFactory          parserFactory
	timeFacade             time.Facade
}

func newObservableRestSourceCreator(
	restRequesterGenerator RestRequesterGenerator,
	parserFactory parserFactory,
	timeFacade time.Facade,
) SourceCreator {
	return &observableRestSourceCreator{
		restRequesterGenerator: restRequesterGenerator,
		parserFactory:          parserFactory,
		timeFacade:             timeFacade,
	}
}

func (creator observableRestSourceCreator) Accept(
	config flam.Bag,
) bool {
	return config.String("driver") == SourceDriverObservableRest &&
		config.Has("uri") &&
		config.Has("path.config") &&
		config.Has("path.timestamp")
}

func (creator observableRestSourceCreator) Create(
	config flam.Bag,
) (Source, error) {
	requester, e := creator.restRequesterGenerator.Create()
	if e != nil {
		return nil, e
	}

	parserId := config.String("parser", DefaultRestParser)
	parser, e := creator.parserFactory.Get(parserId)
	if e != nil {
		return nil, e
	}

	return newObservableRestSource(
		config.Int("priority"),
		requester,
		config.String("uri"),
		parser,
		config.String("path.config"),
		config.String("path.timestamp"),
		creator.timeFacade)
}
