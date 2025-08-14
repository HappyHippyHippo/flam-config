package config

import (
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
)

type sourceFactory flam.Factory[Source]

type sourceFactoryArgs struct {
	dig.In

	Creators      []SourceCreator `group:"flam.config.sources.creator"`
	FactoryConfig flam.FactoryConfig
}

func newSourceFactory(
	args sourceFactoryArgs,
) (sourceFactory, error) {
	var creators []flam.ResourceCreator[Source]
	for _, creator := range args.Creators {
		creators = append(creators, creator)
	}

	return flam.NewFactory(
		creators,
		PathSources,
		args.FactoryConfig,
		nil,
	)
}
