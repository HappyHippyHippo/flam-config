package config

import (
	"io"
	"sync"

	flam "github.com/happyhippyhippo/flam"
)

type Source interface {
	io.Closer

	GetPriority() int
	SetPriority(priority int)
	Get(path string, def ...any) any
}

type ObservableSource interface {
	Source

	Reload() (bool, error)
}

type source struct {
	mutex    sync.Locker
	bag      flam.Bag
	priority int
}

func (*source) Close() error {
	return nil
}

func (source *source) GetPriority() int {
	return source.priority
}

func (source *source) SetPriority(priority int) {
	source.priority = priority
}

func (source *source) Get(
	path string,
	def ...any,
) any {
	source.mutex.Lock()
	defer source.mutex.Unlock()

	return source.bag.Get(path, def...)
}
