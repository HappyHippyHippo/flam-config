package config

import (
	"sync"
	"time"

	flam "github.com/happyhippyhippo/flam"
	flamTime "github.com/happyhippyhippo/flam-time"
)

type observableRestSource struct {
	restSource

	timestampPath string
	timestamp     time.Time
	timeFacade    flamTime.Facade
}

func newObservableRestSource(
	priority int,
	restRequester RestRequester,
	uri string,
	parser Parser,
	configPath string,
	timestampPath string,
	timeFacade flamTime.Facade,
) (Source, error) {
	source := &observableRestSource{
		restSource: restSource{
			source: source{
				mutex:    &sync.Mutex{},
				bag:      flam.Bag{},
				priority: priority,
			},
			uri:           uri,
			configPath:    configPath,
			restRequester: restRequester,
			parser:        parser,
		},
		timestampPath: timestampPath,
		timestamp:     timeFacade.Now(),
		timeFacade:    timeFacade,
	}

	if _, e := source.Reload(); e != nil {
		return nil, e
	}

	return source, nil
}

func (source *observableRestSource) Reload() (bool, error) {
	response, e := source.request()
	if e != nil {
		return false, e
	}

	timestamp, e := source.getTimestamp(response)
	if e != nil {
		return false, e
	}

	if source.timestamp.Equal(source.timeFacade.Unix(0, 0)) || source.timestamp.Before(timestamp) {
		bag, e := source.getConfig(response)
		if e != nil {
			return false, e
		}

		source.mutex.Lock()
		source.bag = bag
		source.timestamp = timestamp
		source.mutex.Unlock()

		return true, nil
	}
	return false, nil
}

func (source *observableRestSource) getTimestamp(
	response flam.Bag,
) (time.Time, error) {
	timestamp := response.Get(source.timestampPath)
	if timestamp == nil {
		return time.Unix(0, 0), newErrRestTimestampNotFound(source.timestampPath, response)
	}

	stringTimestamp, ok := timestamp.(string)
	if !ok {
		return time.Unix(0, 0), newErrRestInvalidTimestamp(source.timestampPath, timestamp)
	}

	parsedTimestamp, e := source.timeFacade.Parse(time.RFC3339, stringTimestamp)
	if e != nil {
		return time.Unix(0, 0), e
	}

	return parsedTimestamp, nil
}
