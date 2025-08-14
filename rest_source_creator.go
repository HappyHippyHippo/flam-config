package config

import (
	flam "github.com/happyhippyhippo/flam"
)

type restSourceCreator struct {
	restRequesterGenerator RestRequesterGenerator
	parserFactory          parserFactory
}

func newRestSourceCreator(
	restRequesterGenerator RestRequesterGenerator,
	parserFactory parserFactory,
) SourceCreator {
	return &restSourceCreator{
		restRequesterGenerator: restRequesterGenerator,
		parserFactory:          parserFactory,
	}
}

func (creator restSourceCreator) Accept(
	config flam.Bag,
) bool {
	return config.String("driver") == SourceDriverRest &&
		config.Has("uri") &&
		config.Has("path.config")
}

func (creator restSourceCreator) Create(
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

	return newRestSource(
		config.Int("priority"),
		requester,
		config.String("uri"),
		parser,
		config.String("path.config"))
}
