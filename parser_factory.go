package config

import (
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
)

type parserFactory flam.Factory[Parser]

type parserFactoryArgs struct {
	dig.In

	Creators      []ParserCreator `group:"flam.config.parsers.creator"`
	FactoryConfig flam.FactoryConfig
}

func newParserFactory(
	args parserFactoryArgs,
) (parserFactory, error) {
	var creators []flam.ResourceCreator[Parser]
	for _, creator := range args.Creators {
		creators = append(creators, creator)
	}

	return flam.NewFactory(
		creators,
		PathParsers,
		args.FactoryConfig,
		nil,
	)
}
