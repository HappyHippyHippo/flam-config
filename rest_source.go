package config

import (
	"net/http"
	"sync"

	flam "github.com/happyhippyhippo/flam"
)

type restSource struct {
	source

	restRequester RestRequester
	uri           string
	parser        Parser
	configPath    string
}

func newRestSource(
	priority int,
	restRequester RestRequester,
	uri string,
	parser Parser,
	configPath string,
) (Source, error) {
	source := &restSource{
		source: source{
			mutex:    &sync.Mutex{},
			bag:      flam.Bag{},
			priority: priority,
		},
		restRequester: restRequester,
		uri:           uri,
		parser:        parser,
		configPath:    configPath,
	}

	if e := source.load(); e != nil {
		return nil, e
	}

	return source, nil
}

func (source *restSource) load() error {
	response, e := source.request()
	if e != nil {
		return e
	}

	bag, e := source.getConfig(response)
	if e != nil {
		return e
	}

	source.mutex.Lock()
	source.bag = bag
	source.mutex.Unlock()

	return nil
}

func (source *restSource) request() (flam.Bag, error) {
	request, e := http.NewRequest(http.MethodGet, source.uri, http.NoBody)
	if e != nil {
		return nil, e
	}

	response, e := source.restRequester.Do(request)
	if e != nil {
		return nil, e
	}

	return source.parser.Parse(response.Body)
}

func (source *restSource) getConfig(
	response flam.Bag,
) (flam.Bag, error) {
	config := response.Get(source.configPath)
	if config == nil {
		return flam.Bag{}, newErrRestConfigNotFound(source.configPath, response)
	}

	bag, ok := config.(flam.Bag)
	if !ok {
		return flam.Bag{}, newErrRestInvalidConfig(source.configPath, config)
	}

	return bag, nil
}
