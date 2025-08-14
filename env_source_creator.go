package config

import (
	flam "github.com/happyhippyhippo/flam"
)

type envSourceCreator struct{}

func newEnvSourceCreator() SourceCreator {
	return &envSourceCreator{}
}

func (envSourceCreator) Accept(
	config flam.Bag,
) bool {
	return config.String("driver") == SourceDriverEnv
}

func (envSourceCreator) Create(
	config flam.Bag,
) (Source, error) {
	mappings := map[string]string{}
	for key, path := range config.Bag("mappings") {
		if str, ok := path.(string); ok {
			mappings[key] = str
		}
	}

	return newEnvSource(
		config.Int("priority"),
		config.StringSlice("files", []string{}),
		mappings,
	)
}
