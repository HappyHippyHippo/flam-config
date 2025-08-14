package config

import (
	flam "github.com/happyhippyhippo/flam"
)

type SourceCreator interface {
	Accept(config flam.Bag) bool
	Create(config flam.Bag) (Source, error)
}
