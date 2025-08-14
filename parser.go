package config

import (
	"io"

	flam "github.com/happyhippyhippo/flam"
)

type Parser interface {
	io.Closer

	Parse(reader io.Reader) (flam.Bag, error)
}
