package config

import (
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"

	flam "github.com/happyhippyhippo/flam"
)

type envSource struct {
	source

	files    []string
	mappings map[string]string
}

func newEnvSource(
	priority int,
	files []string,
	mappings map[string]string,
) (Source, error) {
	source := &envSource{
		source: source{
			mutex:    &sync.Mutex{},
			bag:      flam.Bag{},
			priority: priority,
		},
		files:    files,
		mappings: mappings,
	}

	if e := source.load(); e != nil {
		return nil, e
	}

	return source, nil
}

func (source *envSource) load() error {
	if len(source.files) != 0 {
		if e := godotenv.Load(source.files...); e != nil {
			return e
		}
	}

	for key, path := range source.mappings {
		env := os.Getenv(key)
		if env == "" {
			continue
		}

		step := source.bag
		sections := strings.Split(path, ".")
		for i, section := range sections {
			if i != len(sections)-1 {
				if _, ok := step[section]; !ok {
					step[section] = flam.Bag{}
				}

				step = step[section].(flam.Bag)
			} else {
				step[section] = env
			}
		}
	}

	return nil
}
