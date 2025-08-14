package config

import (
	"encoding/json"
	"io"

	flam "github.com/happyhippyhippo/flam"
)

type jsonParser struct{}

func newJsonParser() Parser {
	return &jsonParser{}
}

func (parser jsonParser) Close() error {
	return nil
}

func (parser jsonParser) Parse(
	reader io.Reader,
) (flam.Bag, error) {
	b, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}

	data := map[string]any{}
	if e := json.Unmarshal(b, &data); e != nil {
		return nil, e
	}

	return Convert(data).(flam.Bag), nil
}
