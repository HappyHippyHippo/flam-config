package config

import (
	"sync"
	"time"

	flam "github.com/happyhippyhippo/flam"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
	flamTime "github.com/happyhippyhippo/flam-time"
)

type observableFileSource struct {
	fileSource

	timeFacade flamTime.Facade
	timestamp  time.Time
}

func newObservableFileSource(
	priority int,
	disk filesystem.Disk,
	path string,
	parser Parser,
	timeFacade flamTime.Facade,
) (Source, error) {
	source := &observableFileSource{
		fileSource: fileSource{
			source: source{
				mutex:    &sync.Mutex{},
				bag:      flam.Bag{},
				priority: priority,
			},
			disk:   disk,
			path:   path,
			parser: parser,
		},
		timeFacade: timeFacade,
		timestamp:  timeFacade.Unix(0, 0),
	}

	if _, e := source.Reload(); e != nil {
		return nil, e
	}

	return source, nil
}

func (source *observableFileSource) Reload() (bool, error) {
	fileStats, e := source.disk.Stat(source.path)
	if e != nil {
		return false, e
	}

	modTime := fileStats.ModTime()
	if source.timestamp.Equal(source.timeFacade.Unix(0, 0)) || source.timestamp.Before(modTime) {
		if e := source.load(); e != nil {
			return false, e
		}
		source.timestamp = modTime

		return true, nil
	}
	return false, nil
}
