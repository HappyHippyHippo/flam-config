package config

import (
	"time"

	flam "github.com/happyhippyhippo/flam"
)

type Facade interface {
	Entries() []string
	Has(path string) bool
	Get(path string, def ...any) any
	Bool(path string, def ...bool) bool
	Int(path string, def ...int) int
	Int8(path string, def ...int8) int8
	Int16(path string, def ...int16) int16
	Int32(path string, def ...int32) int32
	Int64(path string, def ...int64) int64
	Uint(path string, def ...uint) uint
	Uint8(path string, def ...uint8) uint8
	Uint16(path string, def ...uint16) uint16
	Uint32(path string, def ...uint32) uint32
	Uint64(path string, def ...uint64) uint64
	Float32(path string, def ...float32) float32
	Float64(path string, def ...float64) float64
	String(path string, def ...string) string
	StringMap(path string, def ...map[string]any) map[string]any
	StringMapString(path string, def ...map[string]string) map[string]string
	Slice(path string, def ...[]any) []any
	StringSlice(path string, def ...[]string) []string
	Duration(path string, def ...time.Duration) time.Duration
	Bag(path string, def ...flam.Bag) flam.Bag
	Set(path string, value any) error
	Populate(target any, path ...string) error

	HasParser(id string) bool
	ListParsers() []string
	GetParser(id string) (Parser, error)
	AddParser(id string, source Parser) error

	HasSource(id string) bool
	ListSources() []string
	GetSource(id string) (Source, error)
	AddSource(id string, source Source) error
	SetSourcePriority(id string, priority int) error
	RemoveSource(id string) error
	RemoveAllSources() error
	ReloadSources() error

	HasObserver(id, path string) bool
	AddObserver(id, path string, callback Observer) error
	RemoveObserver(id string) error
}

type facade struct {
	parserFactory parserFactory
	manager       *manager
}

func newFacade(
	parserFactory parserFactory,
	manager *manager,
) Facade {
	return &facade{
		parserFactory: parserFactory,
		manager:       manager,
	}
}
func (facade *facade) Entries() []string {
	return facade.manager.aggregate.Entries()
}

func (facade *facade) Has(
	path string,
) bool {
	return facade.manager.aggregate.Has(path)
}

func (facade *facade) Get(
	path string,
	def ...any,
) any {
	return facade.manager.aggregate.Get(path, def...)
}

func (facade *facade) Bool(
	path string,
	def ...bool,
) bool {
	return facade.manager.aggregate.Bool(path, def...)
}

func (facade *facade) Int(
	path string,
	def ...int,
) int {
	return facade.manager.aggregate.Int(path, def...)
}

func (facade *facade) Int8(
	path string,
	def ...int8,
) int8 {
	return facade.manager.aggregate.Int8(path, def...)
}

func (facade *facade) Int16(
	path string,
	def ...int16,
) int16 {
	return facade.manager.aggregate.Int16(path, def...)
}

func (facade *facade) Int32(
	path string,
	def ...int32,
) int32 {
	return facade.manager.aggregate.Int32(path, def...)
}

func (facade *facade) Int64(
	path string,
	def ...int64,
) int64 {
	return facade.manager.aggregate.Int64(path, def...)
}

func (facade *facade) Uint(
	path string,
	def ...uint,
) uint {
	return facade.manager.aggregate.Uint(path, def...)
}

func (facade *facade) Uint8(
	path string,
	def ...uint8,
) uint8 {
	return facade.manager.aggregate.Uint8(path, def...)
}

func (facade *facade) Uint16(
	path string,
	def ...uint16,
) uint16 {
	return facade.manager.aggregate.Uint16(path, def...)
}

func (facade *facade) Uint32(
	path string,
	def ...uint32,
) uint32 {
	return facade.manager.aggregate.Uint32(path, def...)
}

func (facade *facade) Uint64(
	path string,
	def ...uint64,
) uint64 {
	return facade.manager.aggregate.Uint64(path, def...)
}

func (facade *facade) Float32(
	path string,
	def ...float32,
) float32 {
	return facade.manager.aggregate.Float32(path, def...)
}

func (facade *facade) Float64(
	path string,
	def ...float64,
) float64 {
	return facade.manager.aggregate.Float64(path, def...)
}

func (facade *facade) String(
	path string,
	def ...string,
) string {
	return facade.manager.aggregate.String(path, def...)
}

func (facade *facade) StringMap(
	path string,
	def ...map[string]any,
) map[string]any {
	return facade.manager.aggregate.StringMap(path, def...)
}

func (facade *facade) StringMapString(
	path string,
	def ...map[string]string,
) map[string]string {
	return facade.manager.aggregate.StringMapString(path, def...)
}

func (facade *facade) Slice(
	path string,
	def ...[]any,
) []any {
	return facade.manager.aggregate.Slice(path, def...)
}

func (facade *facade) StringSlice(
	path string,
	def ...[]string,
) []string {
	return facade.manager.aggregate.StringSlice(path, def...)
}

func (facade *facade) Duration(
	path string,
	def ...time.Duration,
) time.Duration {
	return facade.manager.aggregate.Duration(path, def...)
}

func (facade *facade) Bag(
	path string,
	def ...flam.Bag,
) flam.Bag {
	return facade.manager.aggregate.Bag(path, def...)
}

func (facade *facade) Set(
	path string,
	value any,
) error {
	return facade.manager.Set(path, value)
}

func (facade *facade) Populate(
	target any,
	path ...string,
) error {
	return facade.manager.aggregate.Populate(target, path...)
}

func (facade *facade) HasParser(
	id string,
) bool {
	return facade.parserFactory.Has(id)
}

func (facade *facade) ListParsers() []string {
	return facade.parserFactory.List()
}

func (facade *facade) GetParser(
	id string,
) (Parser, error) {
	return facade.parserFactory.Get(id)
}

func (facade *facade) AddParser(
	id string,
	source Parser,
) error {
	return facade.parserFactory.Add(id, source)
}

func (facade *facade) HasSource(id string) bool {
	return facade.manager.HasSource(id)
}

func (facade *facade) ListSources() []string {
	return facade.manager.ListSources()
}

func (facade *facade) GetSource(
	id string,
) (Source, error) {
	return facade.manager.GetSource(id)
}

func (facade *facade) AddSource(
	id string,
	source Source,
) error {
	return facade.manager.AddSource(id, source)
}

func (facade *facade) SetSourcePriority(
	id string,
	priority int,
) error {
	return facade.manager.SetSourcePriority(id, priority)
}

func (facade *facade) RemoveSource(
	id string,
) error {
	return facade.manager.RemoveSource(id)
}

func (facade *facade) RemoveAllSources() error {
	return facade.manager.RemoveAllSources()
}

func (facade *facade) ReloadSources() error {
	return facade.manager.ReloadSources()
}

func (facade *facade) HasObserver(
	id,
	path string,
) bool {
	return facade.manager.HasObserver(id, path)
}

func (facade *facade) AddObserver(
	id,
	path string,
	callback Observer,
) error {
	return facade.manager.AddObserver(id, path, callback)
}

func (facade *facade) RemoveObserver(
	id string,
) error {
	return facade.manager.RemoveObserver(id)
}
