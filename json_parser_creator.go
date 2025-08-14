package config

import (
	flam "github.com/happyhippyhippo/flam"
)

type jsonParserCreator struct{}

func newJsonParserCreator() ParserCreator {
	return &jsonParserCreator{}
}

func (jsonParserCreator) Accept(
	config flam.Bag,
) bool {
	return config.String("driver") == ParserDriverJson
}

func (jsonParserCreator) Create(
	_ flam.Bag,
) (Parser, error) {
	return newJsonParser(), nil
}
