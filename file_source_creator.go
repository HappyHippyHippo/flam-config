package config

import (
	flam "github.com/happyhippyhippo/flam"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
)

type fileSourceCreator struct {
	fileSystemFacade filesystem.Facade
	parserFactory    parserFactory
}

func newFileSourceCreator(
	fileSystemFacade filesystem.Facade,
	parserFactory parserFactory,
) SourceCreator {
	return &fileSourceCreator{
		fileSystemFacade: fileSystemFacade,
		parserFactory:    parserFactory,
	}
}

func (fileSourceCreator) Accept(
	config flam.Bag,
) bool {
	return config.String("driver") == SourceDriverFile &&
		config.Has("path")
}

func (creator fileSourceCreator) Create(
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

	return newFileSource(
		config.Int("priority"),
		disk,
		config.String("path"),
		parser)
}
