package config

import (
	"os"
	"sync"

	flam "github.com/happyhippyhippo/flam"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
)

type dirSource struct {
	source

	disk      filesystem.Disk
	path      string
	parser    Parser
	recursive bool
}

func newDirSource(
	priority int,
	disk filesystem.Disk,
	path string,
	parser Parser,
	recursive bool,
) (Source, error) {
	source := &dirSource{
		source: source{
			mutex:    &sync.Mutex{},
			bag:      flam.Bag{},
			priority: priority,
		},
		disk:      disk,
		path:      path,
		parser:    parser,
		recursive: recursive,
	}

	if e := source.load(); e != nil {
		return nil, e
	}

	return source, nil
}

func (source *dirSource) load() error {
	bag, e := source.loadDir(source.path)
	if e != nil {
		return e
	}

	source.mutex.Lock()
	source.bag = bag
	source.mutex.Unlock()

	return nil
}

func (source *dirSource) loadDir(
	path string,
) (flam.Bag, error) {
	dir, e := source.disk.Open(path)
	if e != nil {
		return nil, e
	}
	defer func() { _ = dir.Close() }()

	files, e := dir.Readdir(0)
	if e != nil {
		return nil, e
	}

	loaded := flam.Bag{}
	for _, file := range files {
		if file.IsDir() {
			if source.recursive {
				partial, e := source.loadDir(path + "/" + file.Name())
				if e != nil {
					return nil, e
				}

				loaded.Merge(partial)
			}
		} else {
			partial, e := source.loadFile(path + "/" + file.Name())
			if e != nil {
				return nil, e
			}

			loaded.Merge(partial)
		}
	}

	return loaded, nil
}

func (source *dirSource) loadFile(
	path string,
) (flam.Bag, error) {
	file, e := source.disk.OpenFile(path, os.O_RDONLY, 0o644)
	if e != nil {
		return nil, e
	}
	defer func() { _ = file.Close() }()

	return source.parser.Parse(file)
}
