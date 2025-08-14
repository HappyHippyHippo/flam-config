package config

import (
	"os"
	"sync"

	flam "github.com/happyhippyhippo/flam"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
)

type fileSource struct {
	source

	disk   filesystem.Disk
	path   string
	parser Parser
}

func newFileSource(
	priority int,
	disk filesystem.Disk,
	path string,
	parser Parser,
) (Source, error) {
	source := &fileSource{
		source: source{
			mutex:    &sync.Mutex{},
			bag:      flam.Bag{},
			priority: priority,
		},
		disk:   disk,
		path:   path,
		parser: parser,
	}

	if e := source.load(); e != nil {
		return nil, e
	}

	return source, nil
}

func (source *fileSource) load() error {
	file, e := source.disk.OpenFile(source.path, os.O_RDONLY, 0o644)
	if e != nil {
		return e
	}
	defer func() { _ = file.Close() }()

	bag, e := source.parser.Parse(file)
	if e != nil {
		return e
	}

	source.mutex.Lock()
	defer source.mutex.Unlock()

	source.bag = bag

	return nil
}
