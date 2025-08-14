package config

import (
	"io"

	"gopkg.in/yaml.v3"

	flam "github.com/happyhippyhippo/flam"
)

type yamlParser struct{}

func newYamlParser() Parser {
	return &yamlParser{}
}

func (parser yamlParser) Close() error {
	return nil
}

func (parser yamlParser) Parse(
	reader io.Reader,
) (flam.Bag, error) {
	b, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}

	data := map[string]any{}
	if e := yaml.Unmarshal(b, &data); e != nil {
		return nil, e
	}

	return Convert(data).(flam.Bag), nil
}
