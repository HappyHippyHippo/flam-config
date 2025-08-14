package config

import (
	flam "github.com/happyhippyhippo/flam"
)

type yamlParserCreator struct{}

func newYamlParserCreator() ParserCreator {
	return &yamlParserCreator{}
}

func (yamlParserCreator) Accept(
	config flam.Bag,
) bool {
	return config.String("driver") == ParserDriverYaml
}

func (yamlParserCreator) Create(
	_ flam.Bag,
) (Parser, error) {
	return newYamlParser(), nil
}
