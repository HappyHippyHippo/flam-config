package config

import (
	"io"
	"reflect"
	"slices"
	"sort"
	"strings"
	"sync"

	flam "github.com/happyhippyhippo/flam"
)

type regSource struct {
	id     string
	source Source
}

type regSourceSorter []regSource

func (sorter regSourceSorter) Len() int {
	return len(sorter)
}

func (sorter regSourceSorter) Swap(i, j int) {
	sorter[i], sorter[j] = sorter[j], sorter[i]
}

func (sorter regSourceSorter) Less(i, j int) bool {
	return sorter[i].source.GetPriority() < sorter[j].source.GetPriority()
}

type regObserver struct {
	current   any
	callbacks map[string]Observer
}

type manager struct {
	locker    sync.Locker
	sources   []regSource
	observers map[string]regObserver
	aggregate flam.Bag
	local     flam.Bag
}

func newManager() *manager {
	return &manager{
		locker:    &sync.Mutex{},
		sources:   []regSource{},
		observers: map[string]regObserver{},
		aggregate: flam.Bag{},
		local:     flam.Bag{},
	}
}

func (manager *manager) Set(
	path string,
	value any,
) error {
	if e := manager.local.Set(path, value); e != nil {
		return e
	}

	manager.locker.Lock()
	defer manager.locker.Unlock()

	manager.rebuild()

	return nil
}

func (manager *manager) HasSource(
	id string,
) bool {
	manager.locker.Lock()
	defer manager.locker.Unlock()

	return slices.ContainsFunc(manager.sources, func(s regSource) bool {
		return s.id == id
	})
}

func (manager *manager) ListSources() []string {
	manager.locker.Lock()
	defer manager.locker.Unlock()

	var list []string
	for _, reg := range manager.sources {
		list = append(list, reg.id)
	}

	slices.SortFunc(list, strings.Compare)

	return list
}

func (manager *manager) GetSource(
	id string,
) (Source, error) {
	manager.locker.Lock()
	defer manager.locker.Unlock()

	searcherFunc := func(reg regSource) bool {
		return reg.id == id
	}

	if i := slices.IndexFunc(manager.sources, searcherFunc); i != -1 {
		return manager.sources[i].source, nil
	}

	return nil, newErrSourceNotFound(id)
}

func (manager *manager) AddSource(
	id string,
	source Source,
) error {
	switch {
	case source == nil:
		return newErrNilReference("source")
	case manager.HasSource(id):
		return newErrDuplicateSource(id)
	}

	manager.locker.Lock()
	defer manager.locker.Unlock()

	manager.sources = append(manager.sources, regSource{id, source})
	sort.Sort(regSourceSorter(manager.sources))
	manager.rebuild()

	return nil
}

func (manager *manager) SetSourcePriority(
	id string,
	priority int,
) error {
	manager.locker.Lock()
	defer manager.locker.Unlock()

	searcherFunc := func(reg regSource) bool {
		return reg.id == id
	}

	if i := slices.IndexFunc(manager.sources, searcherFunc); i != -1 {
		manager.sources[i].source.SetPriority(priority)
		sort.Sort(regSourceSorter(manager.sources))
		manager.rebuild()

		return nil
	}

	return newErrSourceNotFound(id)
}

func (manager *manager) RemoveSource(
	id string,
) error {
	manager.locker.Lock()
	defer manager.locker.Unlock()

	searcherFunc := func(reg regSource) bool {
		return reg.id == id
	}

	i := slices.IndexFunc(manager.sources, searcherFunc)
	if i == -1 {
		return newErrSourceNotFound(id)
	}

	if closer, ok := manager.sources[i].source.(io.Closer); ok {
		if e := closer.Close(); e != nil {
			return e
		}
	}

	manager.sources = append(manager.sources[:i], manager.sources[i+1:]...)
	manager.rebuild()

	return nil
}

func (manager *manager) RemoveAllSources() error {
	manager.locker.Lock()
	defer manager.locker.Unlock()

	for _, reg := range manager.sources {
		if source, ok := reg.source.(io.Closer); ok {
			if e := source.Close(); e != nil {
				return e
			}
		}
	}

	manager.sources = []regSource{}
	manager.rebuild()

	return nil
}

func (manager *manager) ReloadSources() error {
	manager.locker.Lock()
	defer manager.locker.Unlock()

	reloaded := false
	for _, ref := range manager.sources {
		if observable, ok := ref.source.(ObservableSource); ok {
			updated, e := observable.Reload()
			if e != nil {
				return e
			}
			reloaded = reloaded || updated
		}
	}

	if reloaded {
		manager.rebuild()
	}

	return nil
}

func (manager *manager) HasObserver(
	id,
	path string,
) bool {
	if reg, ok := manager.observers[path]; ok {
		if _, ok := reg.callbacks[id]; ok {
			return true
		}
	}
	return false
}

func (manager *manager) AddObserver(
	id,
	path string,
	callback Observer,
) error {
	manager.locker.Lock()
	defer manager.locker.Unlock()

	if callback == nil {
		return newErrNilReference("callback")
	}

	if _, ok := manager.observers[path]; !ok {
		manager.observers[path] = regObserver{
			current:   manager.aggregate.Get(path),
			callbacks: map[string]Observer{},
		}
	} else if _, ok := manager.observers[path].callbacks[id]; ok {
		return newErrDuplicateObserver(path, id)
	}

	manager.observers[path].callbacks[id] = callback

	return nil
}

func (manager *manager) RemoveObserver(
	id string,
) error {
	manager.locker.Lock()
	defer manager.locker.Unlock()

	for _, observer := range manager.observers {
		delete(observer.callbacks, id)
	}

	return nil
}

func (manager *manager) rebuild() {
	updated := flam.Bag{}
	for _, ref := range manager.sources {
		updated.Merge(ref.source.Get("", flam.Bag{}).(flam.Bag))
	}

	updated.Merge(manager.local)
	manager.aggregate = updated

	for path, reg := range manager.observers {
		val := manager.aggregate.Get(path, nil)
		if val != nil && !reflect.DeepEqual(reg.current, val) {
			old := reg.current

			manager.observers[path] = regObserver{
				current:   val,
				callbacks: reg.callbacks,
			}

			for _, callback := range manager.observers[path].callbacks {
				callback(old, val)
			}
		}
	}
}
