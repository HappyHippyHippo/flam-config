package config

import (
	flam "github.com/happyhippyhippo/flam"
)

type ParserCreator interface {
	Accept(config flam.Bag) bool
	Create(config flam.Bag) (Parser, error)
}
