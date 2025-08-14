package config

import (
	flam "github.com/happyhippyhippo/flam"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
)

type dirSourceCreator struct {
	fileSystemFacade filesystem.Facade
	parserFactory    parserFactory
}

func newDirSourceCreator(
	fileSystemFacade filesystem.Facade,
	parserFactory parserFactory,
) SourceCreator {
	return &dirSourceCreator{
		fileSystemFacade: fileSystemFacade,
		parserFactory:    parserFactory,
	}
}

func (creator dirSourceCreator) Accept(
	config flam.Bag,
) bool {
	return config.String("driver") == SourceDriverDir &&
		config.Has("path")
}

func (creator dirSourceCreator) Create(
	config flam.Bag,
) (Source, error) {
	diskId := config.String("disk", DefaultFileDisk)
	disk, e := creator.fileSystemFacade.GetDisk(diskId)
	if e != nil {
		return nil, e
	}

	parserId := config.String("parser", DefaultFileParser)
	parser, e := creator.parserFactory.Get(parserId)
	if e != nil {
		return nil, e
	}

	return newDirSource(
		config.Int("priority"),
		disk,
		config.String("path"),
		parser,
		config.Bool("recursive"))
}
