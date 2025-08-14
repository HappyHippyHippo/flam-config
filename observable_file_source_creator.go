package config

import (
	flam "github.com/happyhippyhippo/flam"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
	time "github.com/happyhippyhippo/flam-time"
)

type observableFileSourceCreator struct {
	fileSourceCreator

	timeFacade time.Facade
}

func newObservableFileSourceCreator(
	fileSystemFacade filesystem.Facade,
	parserFactory parserFactory,
	timeFacade time.Facade,
) SourceCreator {
	return &observableFileSourceCreator{
		fileSourceCreator: fileSourceCreator{
			fileSystemFacade: fileSystemFacade,
			parserFactory:    parserFactory,
		},
		timeFacade: timeFacade,
	}
}

func (creator observableFileSourceCreator) Accept(
	config flam.Bag,
) bool {
	return config.String("driver") == SourceDriverObservableFile &&
		config.Has("path")
}

func (creator observableFileSourceCreator) Create(
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

	return newObservableFileSource(
		config.Int("priority"),
		disk,
		config.String("path"),
		parser,
		creator.timeFacade)
}
