package config

import (
	flam "github.com/happyhippyhippo/flam"
)

type factoryConfig struct {
	manager *manager
}

func newFactoryConfig(
	manager *manager,
) flam.FactoryConfig {
	return &factoryConfig{
		manager: manager,
	}
}

func (config factoryConfig) Get(
	path string,
	def ...any,
) flam.Bag {
	data := config.manager.aggregate.Get(path, def...)
	if bag, ok := data.(flam.Bag); ok {
		return bag
	}

	return flam.Bag{}
}
